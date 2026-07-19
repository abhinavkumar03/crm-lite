package dto

import (
	"encoding/json"
	"time"
)

type WorkflowInfo struct {
	Key         string   `json:"workflow_key"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     int      `json:"version"`
	DurationMin int      `json:"duration_min"`
	Features    []string `json:"features"`
}

type WorkflowDefinition struct {
	Workflow WorkflowInfo `json:"workflow"`
	Steps    []StepDTO    `json:"steps"`
}

type StepDTO struct {
	Key             string          `json:"step_key"`
	SortOrder       int             `json:"sort_order"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	WhyItMatters    string          `json:"why_it_matters"`
	HowItWorks      string          `json:"how_it_works"`
	ExpectedResult  string          `json:"expected_result"`
	Route           *string         `json:"route,omitempty"`
	TargetSelector  *string         `json:"target_selector,omitempty"`
	ActionLabel     *string         `json:"action_label,omitempty"`
	ValidatorKey    string          `json:"validator_key"`
	ValidatorParams json.RawMessage `json:"validator_params"`
	IsSkippable     bool            `json:"is_skippable"`
	Status          string          `json:"status"`
	RequiredAction  string          `json:"required_action,omitempty"`
	SuccessEvent    *string         `json:"success_event,omitempty"`
	FailureMessage  string          `json:"failure_message,omitempty"`
	Hint            string          `json:"hint,omitempty"`
	MaxAttempts     int             `json:"max_attempts,omitempty"`
	AllowSelectors  json.RawMessage `json:"allow_selectors,omitempty"`
	Placement       string          `json:"placement,omitempty"`
}

type SessionDTO struct {
	ID                    string          `json:"id"`
	WorkflowKey           string          `json:"workflow_key"`
	WorkflowVersion       int             `json:"workflow_version"`
	SandboxOrganizationID *string         `json:"sandbox_organization_id,omitempty"`
	Status                string          `json:"status"`
	CurrentStepKey        *string         `json:"current_step_key,omitempty"`
	StartedAt             time.Time       `json:"started_at"`
	CompletedAt           *time.Time      `json:"completed_at,omitempty"`
	Stats                 json.RawMessage `json:"stats"`
	Steps                 []StepDTO       `json:"steps"`
	ProgressPercent       int             `json:"progress_percent"`
}

type ValidateResult struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message"`
	Session *SessionDTO `json:"session,omitempty"`
}

type CleanupRequest struct {
	KeepData bool `json:"keep_data"`
}

type ClientEvent struct {
	Type     string `json:"type"`
	Selector string `json:"selector"`
	Path     string `json:"path"`
}

type ValidateRequest struct {
	StepKey     string       `json:"step_key"`
	Route       string       `json:"route"`
	ClientEvent *ClientEvent `json:"client_event"`
}

type EventRequest struct {
	EventType string          `json:"event_type" validate:"required"`
	Payload   json.RawMessage `json:"payload"`
}
