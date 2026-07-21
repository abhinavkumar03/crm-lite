package dto

import "time"

// TriggerInput defines a trigger on create/update.
type TriggerInput struct {
	Type   string         `json:"type" binding:"required"`
	Config map[string]any `json:"config"`
}

// ConditionInput is a nested condition tree node.
type ConditionInput struct {
	NodeType     string           `json:"node_type" binding:"required"`
	Logic        *string          `json:"logic,omitempty"`
	FieldAPIName *string          `json:"field_api_name,omitempty"`
	Operator     *string          `json:"operator,omitempty"`
	Value        any              `json:"value,omitempty"`
	Children     []ConditionInput `json:"children,omitempty"`
}

// ActionInput defines an ordered action step.
type ActionInput struct {
	Type            string         `json:"type" binding:"required"`
	Config          map[string]any `json:"config"`
	MaxRetries      int            `json:"max_retries"`
	ContinueOnError *bool          `json:"continue_on_error,omitempty"`
}

// CreateWorkflowRequest creates a draft workflow with optional definition.
type CreateWorkflowRequest struct {
	Name          string           `json:"name" binding:"required,min=1,max=150"`
	Description   string           `json:"description"`
	ModuleID      *string          `json:"module_id"`
	OnActionError string           `json:"on_action_error"`
	Priority      int              `json:"priority"`
	Triggers      []TriggerInput   `json:"triggers"`
	Conditions    *ConditionInput  `json:"conditions"`
	Actions       []ActionInput    `json:"actions"`
}

// UpdateWorkflowRequest patches header and/or draft definition.
type UpdateWorkflowRequest struct {
	Name          *string          `json:"name"`
	Description   *string          `json:"description"`
	ModuleID      *string          `json:"module_id"`
	OnActionError *string          `json:"on_action_error"`
	Priority      *int             `json:"priority"`
	Triggers      []TriggerInput   `json:"triggers"`
	Conditions    *ConditionInput  `json:"conditions"`
	Actions       []ActionInput    `json:"actions"`
}

// PublishRequest optional changelog on publish.
type PublishRequest struct {
	Changelog string `json:"changelog"`
}

// ManualRunRequest runs a workflow against a record.
type ManualRunRequest struct {
	RecordID string `json:"record_id" binding:"required"`
	ModuleID string `json:"module_id"`
}

// TriggerResponse is an API trigger.
type TriggerResponse struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

// ConditionResponse is a nested condition tree.
type ConditionResponse struct {
	ID           string              `json:"id"`
	NodeType     string              `json:"node_type"`
	Logic        *string             `json:"logic,omitempty"`
	FieldAPIName *string             `json:"field_api_name,omitempty"`
	Operator     *string             `json:"operator,omitempty"`
	Value        any                 `json:"value,omitempty"`
	Children     []ConditionResponse `json:"children,omitempty"`
}

// ActionResponse is an ordered action.
type ActionResponse struct {
	ID              string         `json:"id"`
	SortOrder       int            `json:"sort_order"`
	Type            string         `json:"type"`
	Config          map[string]any `json:"config"`
	MaxRetries      int            `json:"max_retries"`
	ContinueOnError *bool          `json:"continue_on_error,omitempty"`
}

// WorkflowResponse is the full workflow + current draft/published definition.
type WorkflowResponse struct {
	ID                 string              `json:"id"`
	ModuleID           *string             `json:"module_id"`
	ModuleAPIName      *string             `json:"module_api_name,omitempty"`
	Name               string              `json:"name"`
	Description        string              `json:"description"`
	Status             string              `json:"status"`
	OnActionError      string              `json:"on_action_error"`
	Priority           int                 `json:"priority"`
	PublishedVersionID *string             `json:"published_version_id,omitempty"`
	DraftVersionID     *string             `json:"draft_version_id,omitempty"`
	Version            int                 `json:"version"`
	Triggers           []TriggerResponse   `json:"triggers"`
	Conditions         *ConditionResponse  `json:"conditions,omitempty"`
	Actions            []ActionResponse    `json:"actions"`
	CreatedBy          *string             `json:"created_by,omitempty"`
	UpdatedBy          *string             `json:"updated_by,omitempty"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
}

// WorkflowSummary is a list row.
type WorkflowSummary struct {
	ID            string    `json:"id"`
	ModuleID      *string   `json:"module_id"`
	ModuleAPIName *string   `json:"module_api_name,omitempty"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Status        string    `json:"status"`
	Priority      int       `json:"priority"`
	Version       int       `json:"version"`
	TriggerTypes  []string  `json:"trigger_types"`
	ActionCount   int       `json:"action_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ListWorkflowsResult paginated workflows.
type ListWorkflowsResult struct {
	Items      []WorkflowSummary `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	Total      int               `json:"total"`
	TotalPages int               `json:"total_pages"`
}

// VersionSummary lists version history.
type VersionSummary struct {
	ID          string     `json:"id"`
	Version     int        `json:"version"`
	State       string     `json:"state"`
	Changelog   string     `json:"changelog"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	PublishedBy *string    `json:"published_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// ExecutionSummary list row.
type ExecutionSummary struct {
	ID           string     `json:"id"`
	WorkflowID   string     `json:"workflow_id"`
	WorkflowName string     `json:"workflow_name,omitempty"`
	VersionID    *string    `json:"version_id,omitempty"`
	ModuleID     *string    `json:"module_id,omitempty"`
	RecordID     *string    `json:"record_id,omitempty"`
	TriggerType  string     `json:"trigger_type"`
	Status       string     `json:"status"`
	Source       string     `json:"source"`
	Depth        int        `json:"depth"`
	ErrorSummary *string    `json:"error_summary,omitempty"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	DurationMs   *int       `json:"duration_ms,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

// ExecutionStepResponse one step.
type ExecutionStepResponse struct {
	ID         string         `json:"id"`
	SortOrder  int            `json:"sort_order"`
	ActionType string         `json:"action_type"`
	ActionID   *string        `json:"action_id,omitempty"`
	Status     string         `json:"status"`
	Input      map[string]any `json:"input"`
	Output     map[string]any `json:"output"`
	Error      *string        `json:"error,omitempty"`
	StartedAt  *time.Time     `json:"started_at,omitempty"`
	FinishedAt *time.Time     `json:"finished_at,omitempty"`
}

// ExecutionDetail full run with steps.
type ExecutionDetail struct {
	ExecutionSummary
	Steps []ExecutionStepResponse `json:"steps"`
}

// ListExecutionsResult paginated executions.
type ListExecutionsResult struct {
	Items      []ExecutionSummary `json:"items"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	Total      int                `json:"total"`
	TotalPages int                `json:"total_pages"`
}

// TemplateResponse starter template.
type TemplateResponse struct {
	ID            string         `json:"id"`
	Key           string         `json:"key"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	ModuleAPIName *string        `json:"module_api_name,omitempty"`
	Category      string         `json:"category,omitempty"`
	Definition    map[string]any `json:"definition"`
	IsBuiltin     bool           `json:"is_builtin"`
}

// BuilderMetadataResponse options for the visual builder.
type BuilderMetadataResponse struct {
	Modules   []BuilderModule   `json:"modules"`
	Operators []BuilderOperator `json:"operators"`
	Actions   []BuilderAction   `json:"actions"`
	Triggers  []BuilderTrigger  `json:"triggers"`
	Variables []BuilderVariable `json:"variables"`
	Users     []BuilderUser     `json:"users"`
	Templates []TemplateResponse `json:"templates"`
}

type BuilderModule struct {
	ID      string         `json:"id"`
	APIName string         `json:"api_name"`
	Label   string         `json:"label"`
	Fields  []BuilderField `json:"fields"`
}

type BuilderField struct {
	APIName  string   `json:"api_name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Options  []string `json:"options,omitempty"`
	Required bool     `json:"required"`
}

type BuilderOperator struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	ValueArity  string   `json:"value_arity"` // none | one | two | list
	FieldTypes  []string `json:"field_types,omitempty"`
}

type BuilderAction struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description"`
	MVP         bool   `json:"mvp"`
}

type BuilderTrigger struct {
	Type        string `json:"type"`
	Label       string `json:"label"`
	Description string `json:"description"`
	MVP         bool   `json:"mvp"`
}

type BuilderVariable struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type BuilderUser struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// MetricsResponse dashboard rollups.
type MetricsResponse struct {
	ActiveWorkflows   int      `json:"active_workflows"`
	DisabledWorkflows int      `json:"disabled_workflows"`
	DraftWorkflows    int      `json:"draft_workflows"`
	ExecutedToday     int      `json:"executed_today"`
	FailedToday       int      `json:"failed_today"`
	AvgDurationMs     *float64 `json:"avg_duration_ms,omitempty"`
}
