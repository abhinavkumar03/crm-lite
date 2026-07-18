package dto

// ValidateRequest is a payload to dry-run validation against a module's schema.
type ValidateRequest struct {
	Data map[string]any `json:"data"`
}

// FieldError is a single validation failure, keyed by the field's api_name.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateResult is the outcome of running the validator.
type ValidateResult struct {
	Valid  bool         `json:"valid"`
	Errors []FieldError `json:"errors"`
}

// FieldSchema is the compiled, frontend-consumable constraint set for a field.
// Pointers are omitted when not applicable so the client only sees active rules.
type FieldSchema struct {
	APIName   string            `json:"api_name"`
	Label     string            `json:"label"`
	Type      string            `json:"type"`
	Required  bool              `json:"required"`
	MinLength *int              `json:"min_length,omitempty"`
	MaxLength *int              `json:"max_length,omitempty"`
	Min       *float64          `json:"min,omitempty"`
	Max       *float64          `json:"max,omitempty"`
	Pattern   *string           `json:"pattern,omitempty"`
	Format    *string           `json:"format,omitempty"` // email | url
	Options   []string          `json:"options,omitempty"`
	Multiple  bool              `json:"multiple,omitempty"`
	Messages  map[string]string `json:"messages,omitempty"` // ruleType -> custom message
}

// ValidationSchema is the full compiled schema for a module.
type ValidationSchema struct {
	ModuleID string        `json:"module_id"`
	Fields   []FieldSchema `json:"fields"`
}
