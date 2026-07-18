package dto

import (
	"encoding/json"
	"time"
)

// RuleResponse is the API representation of a validation rule.
type RuleResponse struct {
	ID           string          `json:"id"`
	ModuleID     string          `json:"module_id"`
	FieldID      *string         `json:"field_id"`
	RuleType     string          `json:"rule_type"`
	Params       json.RawMessage `json:"params"`
	ErrorMessage *string         `json:"error_message"`
	IsActive     bool            `json:"is_active"`
	SortOrder    int             `json:"sort_order"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// CreateRuleRequest is the payload for adding a validation rule.
type CreateRuleRequest struct {
	FieldID      *string        `json:"field_id" validate:"omitempty,uuid"`
	RuleType     string         `json:"rule_type" validate:"required"`
	Params       map[string]any `json:"params"`
	ErrorMessage *string        `json:"error_message" validate:"omitempty,max=300"`
	IsActive     *bool          `json:"is_active"`
	SortOrder    *int           `json:"sort_order"`
}

// UpdateRuleRequest is a partial update. rule_type / field_id are immutable.
type UpdateRuleRequest struct {
	Params       map[string]any `json:"params"`
	ErrorMessage *string        `json:"error_message" validate:"omitempty,max=300"`
	IsActive     *bool          `json:"is_active"`
	SortOrder    *int           `json:"sort_order"`
}
