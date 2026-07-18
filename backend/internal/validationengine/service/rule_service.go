package service

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"

	fieldentity "github.com/abhinavkumar03/crm-lite/backend/internal/field/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/entity"
)

// Domain errors mapped to HTTP status codes by the handler.
var (
	ErrModuleNotFound  = errors.New("module not found")
	ErrNotFound        = errors.New("validation rule not found")
	ErrInvalidRuleType = errors.New("unsupported rule_type")
	ErrFieldRequired   = errors.New("this rule_type requires a valid field_id belonging to the module")
	ErrModuleRule      = errors.New("this rule_type is module-level and must not target a field")
	ErrInvalidParams   = errors.New("invalid params for this rule_type")
)

// RuleRepository is the persistence contract for validation rules.
type RuleRepository interface {
	ModuleExists(ctx context.Context, orgID, moduleID string) (bool, error)
	FieldExists(ctx context.Context, orgID, moduleID, fieldID string) (bool, error)
	Create(ctx context.Context, rule *entity.ValidationRule) error
	List(ctx context.Context, orgID, moduleID string) ([]entity.ValidationRule, error)
	ActiveByModule(ctx context.Context, orgID, moduleID string) ([]entity.ValidationRule, error)
	GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.ValidationRule, error)
	Update(ctx context.Context, rule *entity.ValidationRule) error
	Delete(ctx context.Context, orgID, moduleID, id string) (bool, error)
}

// FieldReader lets the engine read a module's field metadata. It is satisfied by
// the field engine's repository (dependency inversion, no import cycle).
type FieldReader interface {
	ModuleStorage(ctx context.Context, orgID, moduleID string) (string, bool, error)
	List(ctx context.Context, orgID, moduleID string) ([]fieldentity.Field, error)
}

type Service struct {
	rules  RuleRepository
	fields FieldReader
}

func New(rules RuleRepository, fields FieldReader) *Service {
	return &Service{rules: rules, fields: fields}
}

func (s *Service) List(ctx context.Context, orgID, moduleID string) ([]dto.RuleResponse, error) {
	ok, err := s.rules.ModuleExists(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	rules, err := s.rules.List(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.RuleResponse, 0, len(rules))
	for i := range rules {
		out = append(out, toRuleResponse(&rules[i]))
	}
	return out, nil
}

func (s *Service) GetByID(ctx context.Context, orgID, moduleID, id string) (*dto.RuleResponse, error) {
	rule, err := s.rules.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, nil
	}
	resp := toRuleResponse(rule)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, orgID, moduleID string, req dto.CreateRuleRequest) (*dto.RuleResponse, error) {
	ok, err := s.rules.ModuleExists(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrModuleNotFound
	}

	fieldLevel := entity.FieldLevelTypes[req.RuleType]
	moduleLevel := entity.ModuleLevelTypes[req.RuleType]
	if !fieldLevel && !moduleLevel {
		return nil, ErrInvalidRuleType
	}

	if fieldLevel {
		if req.FieldID == nil || *req.FieldID == "" {
			return nil, ErrFieldRequired
		}
		exists, err := s.rules.FieldExists(ctx, orgID, moduleID, *req.FieldID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, ErrFieldRequired
		}
	} else if req.FieldID != nil && *req.FieldID != "" {
		return nil, ErrModuleRule
	}

	if err := validateParams(req.RuleType, req.Params); err != nil {
		return nil, err
	}

	params, err := marshalParams(req.Params)
	if err != nil {
		return nil, err
	}

	var fieldID *string
	if fieldLevel {
		fieldID = req.FieldID
	}

	rule := &entity.ValidationRule{
		OrganizationID: orgID,
		ModuleID:       moduleID,
		FieldID:        fieldID,
		RuleType:       req.RuleType,
		Params:         params,
		ErrorMessage:   req.ErrorMessage,
		IsActive:       derefBool(req.IsActive, true),
		SortOrder:      derefInt(req.SortOrder, 0),
	}

	if err := s.rules.Create(ctx, rule); err != nil {
		return nil, err
	}

	resp := toRuleResponse(rule)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, orgID, moduleID, id string, req dto.UpdateRuleRequest) (*dto.RuleResponse, error) {
	rule, err := s.rules.GetByID(ctx, orgID, moduleID, id)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, ErrNotFound
	}

	if req.Params != nil {
		if err := validateParams(rule.RuleType, req.Params); err != nil {
			return nil, err
		}
		params, err := marshalParams(req.Params)
		if err != nil {
			return nil, err
		}
		rule.Params = params
	}
	if req.ErrorMessage != nil {
		rule.ErrorMessage = req.ErrorMessage
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		rule.SortOrder = *req.SortOrder
	}

	if err := s.rules.Update(ctx, rule); err != nil {
		return nil, err
	}

	resp := toRuleResponse(rule)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, orgID, moduleID, id string) error {
	deleted, err := s.rules.Delete(ctx, orgID, moduleID, id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

// validateParams ensures the params map is well-formed for the rule type.
func validateParams(ruleType string, params map[string]any) error {
	switch ruleType {
	case entity.RuleRequired, entity.RuleEmail, entity.RuleURL:
		return nil
	case entity.RuleMinLength, entity.RuleMaxLength, entity.RuleMin, entity.RuleMax:
		if _, ok := paramFloat(params, "value"); !ok {
			return ErrInvalidParams
		}
		return nil
	case entity.RulePattern:
		pat, ok := paramString(params, "pattern")
		if !ok || pat == "" {
			return ErrInvalidParams
		}
		if _, err := regexp.Compile(pat); err != nil {
			return ErrInvalidParams
		}
		return nil
	case entity.RuleIn, entity.RuleNotIn:
		if len(paramStrings(params, "values")) == 0 {
			return ErrInvalidParams
		}
		return nil
	case entity.RuleRequiredIf:
		field, okF := paramString(params, "field")
		target, okT := paramString(params, "target")
		_, hasEquals := params["equals"]
		if !okF || field == "" || !okT || target == "" || !hasEquals {
			return ErrInvalidParams
		}
		return nil
	default:
		return ErrInvalidRuleType
	}
}

func marshalParams(params map[string]any) ([]byte, error) {
	if len(params) == 0 {
		return []byte("{}"), nil
	}
	return json.Marshal(params)
}

func toRuleResponse(r *entity.ValidationRule) dto.RuleResponse {
	params := json.RawMessage(r.Params)
	if len(params) == 0 {
		params = json.RawMessage("{}")
	}
	return dto.RuleResponse{
		ID:           r.ID,
		ModuleID:     r.ModuleID,
		FieldID:      r.FieldID,
		RuleType:     r.RuleType,
		Params:       params,
		ErrorMessage: r.ErrorMessage,
		IsActive:     r.IsActive,
		SortOrder:    r.SortOrder,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func derefBool(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

func derefInt(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}
