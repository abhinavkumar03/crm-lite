package service

import (
	"context"
	"regexp"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/entity"
)

// Validate evaluates a data payload against a module's fields and active,
// database-driven validation rules. It returns structured, field-keyed errors
// with custom messages. This engine is reused by the record runtime (Phase 10)
// and exposed directly via a dry-run endpoint.
func (s *Service) Validate(ctx context.Context, orgID, moduleID string, data map[string]any) (dto.ValidateResult, error) {
	fields, rules, err := s.load(ctx, orgID, moduleID)
	if err != nil {
		return dto.ValidateResult{}, err
	}

	byField, moduleRules := groupRules(rules)

	errsList := make([]dto.FieldError, 0)
	addErr := func(field, message string) {
		errsList = append(errsList, dto.FieldError{Field: field, Message: message})
	}

	for i := range fields {
		f := fields[i]
		value := data[f.APIName]
		s.validateField(f, value, byField[f.ID], addErr)
	}

	for i := range moduleRules {
		evalModuleRule(moduleRules[i], data, addErr)
	}

	return dto.ValidateResult{Valid: len(errsList) == 0, Errors: errsList}, nil
}

// Schema compiles a frontend-consumable constraint set per field, merging field
// metadata with database-driven rules.
func (s *Service) Schema(ctx context.Context, orgID, moduleID string) (dto.ValidationSchema, error) {
	fields, rules, err := s.load(ctx, orgID, moduleID)
	if err != nil {
		return dto.ValidationSchema{}, err
	}

	byField, _ := groupRules(rules)

	out := dto.ValidationSchema{ModuleID: moduleID, Fields: make([]dto.FieldSchema, 0, len(fields))}
	for i := range fields {
		out.Fields = append(out.Fields, compileFieldSchema(fields[i], byField[fields[i].ID]))
	}
	return out, nil
}

func (s *Service) load(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, []entity.ValidationRule, error) {
	_, ok, err := s.fields.ModuleStorage(ctx, orgID, moduleID)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, ErrModuleNotFound
	}

	fields, err := s.fields.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, nil, err
	}
	rules, err := s.rules.ActiveByModule(ctx, orgID, moduleID)
	if err != nil {
		return nil, nil, err
	}
	return fields, rules, nil
}

func groupRules(rules []entity.ValidationRule) (map[string][]entity.ValidationRule, []entity.ValidationRule) {
	byField := map[string][]entity.ValidationRule{}
	var moduleRules []entity.ValidationRule
	for _, r := range rules {
		if r.FieldID == nil || *r.FieldID == "" {
			moduleRules = append(moduleRules, r)
			continue
		}
		byField[*r.FieldID] = append(byField[*r.FieldID], r)
	}
	return byField, moduleRules
}

func (s *Service) validateField(f fieldentity.Field, value any, rules []entity.ValidationRule, addErr func(field, message string)) {
	required := f.IsRequired
	var requiredMsg *string
	for _, r := range rules {
		if r.RuleType == entity.RuleRequired {
			required = true
			requiredMsg = r.ErrorMessage
		}
	}

	if isEmpty(value) {
		if required {
			addErr(f.APIName, pick([]*string{requiredMsg, f.ValidationMessage}, "This field is required"))
		}
		return
	}

	validateMetadata(f, value, addErr)

	for _, r := range rules {
		if r.RuleType == entity.RuleRequired {
			continue
		}
		if !evalFieldRule(r.RuleType, value, parseParams(r.Params)) {
			addErr(f.APIName, pick([]*string{r.ErrorMessage}, defaultMessage(r.RuleType)))
		}
	}
}

// validateMetadata applies the constraints implied by the field definition
// itself (type, length, regex, options).
func validateMetadata(f fieldentity.Field, value any, addErr func(field, message string)) {
	api := f.APIName

	switch f.FieldType {
	case fieldentity.TypeEmail:
		if !emailRegex.MatchString(toString(value)) {
			addErr(api, pick([]*string{f.ValidationMessage}, "Must be a valid email"))
		}
	case fieldentity.TypeURL:
		if !isURL(toString(value)) {
			addErr(api, pick([]*string{f.ValidationMessage}, "Must be a valid URL"))
		}
	case fieldentity.TypeNumber, fieldentity.TypeCurrency:
		if _, ok := toFloat(value); !ok {
			addErr(api, pick([]*string{f.ValidationMessage}, "Must be a number"))
		}
	case fieldentity.TypeDropdown, fieldentity.TypeRadio:
		if opts := parseOptionValues(f.Options); len(opts) > 0 && !contains(opts, toString(value)) {
			addErr(api, pick([]*string{f.ValidationMessage}, "Invalid option selected"))
		}
	case fieldentity.TypeMultiselect:
		if opts := parseOptionValues(f.Options); len(opts) > 0 {
			if arr, ok := value.([]any); ok {
				for _, item := range arr {
					if !contains(opts, toString(item)) {
						addErr(api, pick([]*string{f.ValidationMessage}, "Invalid option selected"))
						break
					}
				}
			}
		}
	}

	if stringTypes[f.FieldType] {
		length := runeLen(toString(value))
		if f.MinLength != nil && length < *f.MinLength {
			addErr(api, pick([]*string{f.ValidationMessage}, "Value is too short"))
		}
		if f.MaxLength != nil && length > *f.MaxLength {
			addErr(api, pick([]*string{f.ValidationMessage}, "Value is too long"))
		}
		if f.Regex != nil && *f.Regex != "" {
			if re, err := regexp.Compile(*f.Regex); err == nil && !re.MatchString(toString(value)) {
				addErr(api, pick([]*string{f.ValidationMessage}, "Invalid format"))
			}
		}
	}
}

// evalFieldRule returns true when the value satisfies the rule. Malformed params
// (already rejected at write time) are treated as satisfied.
func evalFieldRule(ruleType string, value any, params map[string]any) bool {
	switch ruleType {
	case entity.RuleMinLength:
		v, ok := paramFloat(params, "value")
		return !ok || runeLen(toString(value)) >= int(v)
	case entity.RuleMaxLength:
		v, ok := paramFloat(params, "value")
		return !ok || runeLen(toString(value)) <= int(v)
	case entity.RuleMin:
		v, ok := paramFloat(params, "value")
		if !ok {
			return true
		}
		n, ok := toFloat(value)
		return !ok || n >= v
	case entity.RuleMax:
		v, ok := paramFloat(params, "value")
		if !ok {
			return true
		}
		n, ok := toFloat(value)
		return !ok || n <= v
	case entity.RulePattern:
		pat, ok := paramString(params, "pattern")
		if !ok {
			return true
		}
		re, err := regexp.Compile(pat)
		return err != nil || re.MatchString(toString(value))
	case entity.RuleEmail:
		return emailRegex.MatchString(toString(value))
	case entity.RuleURL:
		return isURL(toString(value))
	case entity.RuleIn:
		return contains(paramStrings(params, "values"), toString(value))
	case entity.RuleNotIn:
		return !contains(paramStrings(params, "values"), toString(value))
	default:
		return true
	}
}

// evalModuleRule evaluates a cross-field (module-level) rule against the whole
// payload.
func evalModuleRule(r entity.ValidationRule, data map[string]any, addErr func(field, message string)) {
	params := parseParams(r.Params)
	switch r.RuleType {
	case entity.RuleRequiredIf:
		field, _ := paramString(params, "field")
		target, _ := paramString(params, "target")
		equals := params["equals"]
		if toString(data[field]) == toString(equals) && isEmpty(data[target]) {
			addErr(target, pick([]*string{r.ErrorMessage}, "This field is required"))
		}
	}
}

func defaultMessage(ruleType string) string {
	switch ruleType {
	case entity.RuleMinLength:
		return "Value is too short"
	case entity.RuleMaxLength:
		return "Value is too long"
	case entity.RuleMin:
		return "Value is too small"
	case entity.RuleMax:
		return "Value is too large"
	case entity.RulePattern:
		return "Invalid format"
	case entity.RuleEmail:
		return "Must be a valid email"
	case entity.RuleURL:
		return "Must be a valid URL"
	case entity.RuleIn:
		return "Value is not allowed"
	case entity.RuleNotIn:
		return "Value is not allowed"
	default:
		return "Invalid value"
	}
}

// compileFieldSchema merges field metadata and rules into a client-side schema.
func compileFieldSchema(f fieldentity.Field, rules []entity.ValidationRule) dto.FieldSchema {
	fs := dto.FieldSchema{
		APIName:  f.APIName,
		Label:    f.Label,
		Type:     f.FieldType,
		Required: f.IsRequired,
	}

	if stringTypes[f.FieldType] {
		fs.MinLength = f.MinLength
		fs.MaxLength = f.MaxLength
		if f.Regex != nil && *f.Regex != "" {
			fs.Pattern = f.Regex
		}
	}

	switch f.FieldType {
	case fieldentity.TypeEmail:
		fs.Format = strPtr("email")
	case fieldentity.TypeURL:
		fs.Format = strPtr("url")
	case fieldentity.TypeDropdown, fieldentity.TypeRadio:
		fs.Options = parseOptionValues(f.Options)
	case fieldentity.TypeMultiselect:
		fs.Options = parseOptionValues(f.Options)
		fs.Multiple = true
	}

	messages := map[string]string{}
	for _, r := range rules {
		params := parseParams(r.Params)
		switch r.RuleType {
		case entity.RuleRequired:
			fs.Required = true
		case entity.RuleMinLength:
			if v, ok := paramFloat(params, "value"); ok {
				fs.MinLength = intPtr(int(v))
			}
		case entity.RuleMaxLength:
			if v, ok := paramFloat(params, "value"); ok {
				fs.MaxLength = intPtr(int(v))
			}
		case entity.RuleMin:
			if v, ok := paramFloat(params, "value"); ok {
				fs.Min = floatPtr(v)
			}
		case entity.RuleMax:
			if v, ok := paramFloat(params, "value"); ok {
				fs.Max = floatPtr(v)
			}
		case entity.RulePattern:
			if pat, ok := paramString(params, "pattern"); ok {
				fs.Pattern = strPtr(pat)
			}
		case entity.RuleEmail:
			fs.Format = strPtr("email")
		case entity.RuleURL:
			fs.Format = strPtr("url")
		case entity.RuleIn:
			fs.Options = paramStrings(params, "values")
		}
		if r.ErrorMessage != nil && *r.ErrorMessage != "" {
			messages[r.RuleType] = *r.ErrorMessage
		}
	}

	if len(messages) > 0 {
		fs.Messages = messages
	}
	return fs
}

func strPtr(s string) *string     { return &s }
func intPtr(i int) *int           { return &i }
func floatPtr(f float64) *float64 { return &f }
