package entity

import "time"

// Rule types supported by the validation engine. Field-level rules target a
// single field (field_id set); RuleRequiredIf is a cross-field (module-level)
// rule (field_id NULL).
const (
	RuleRequired   = "required"
	RuleMinLength  = "min_length"
	RuleMaxLength  = "max_length"
	RuleMin        = "min"
	RuleMax        = "max"
	RulePattern    = "pattern"
	RuleEmail      = "email"
	RuleURL        = "url"
	RuleIn         = "in"
	RuleNotIn      = "not_in"
	RuleRequiredIf = "required_if"
)

// FieldLevelTypes are rule types that must target a specific field.
var FieldLevelTypes = map[string]bool{
	RuleRequired:  true,
	RuleMinLength: true,
	RuleMaxLength: true,
	RuleMin:       true,
	RuleMax:       true,
	RulePattern:   true,
	RuleEmail:     true,
	RuleURL:       true,
	RuleIn:        true,
	RuleNotIn:     true,
}

// ModuleLevelTypes are rule types that operate across fields (field_id NULL).
var ModuleLevelTypes = map[string]bool{
	RuleRequiredIf: true,
}

// ValidationRule is a database-driven rule attached to a module or a field.
type ValidationRule struct {
	ID             string
	OrganizationID string
	ModuleID       string
	FieldID        *string
	RuleType       string
	Params         []byte // raw JSONB
	ErrorMessage   *string
	IsActive       bool
	SortOrder      int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
