package entity

import (
	"encoding/json"
	"time"
)

// Workflow lifecycle statuses.
const (
	StatusDraft    = "draft"
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusArchived = "archived"
)

const (
	OnErrorContinue = "continue"
	OnErrorStop     = "stop"
)

const (
	VersionDraft      = "draft"
	VersionPublished  = "published"
	VersionRolledBack = "rolled_back"
)

const (
	TriggerRecordCreated = "record_created"
	TriggerRecordUpdated = "record_updated"
	TriggerRecordDeleted = "record_deleted"
	TriggerFieldUpdated  = "field_updated"
	TriggerScheduled     = "scheduled"
	TriggerDateBased     = "date_based"
	TriggerManual        = "manual"
)

const (
	NodeGroup     = "group"
	NodePredicate = "predicate"
	LogicAnd      = "and"
	LogicOr       = "or"
)

const (
	ActionUpdateRecord   = "update_record"
	ActionCreateRecord   = "create_record"
	ActionDeleteRecord   = "delete_record"
	ActionAssignOwner    = "assign_owner"
	ActionSendEmail      = "send_email"
	ActionSendWhatsApp   = "send_whatsapp"
	ActionCreateNote     = "create_note"
	ActionCreateActivity = "create_activity"
	ActionWebhook        = "webhook"
	ActionDelay          = "delay"
	ActionInvokeWorkflow = "invoke_workflow"
	ActionBranch         = "branch"
)

const (
	ExecQueued     = "queued"
	ExecRunning    = "running"
	ExecSucceeded  = "succeeded"
	ExecFailed     = "failed"
	ExecPartial    = "partial"
	ExecCancelled  = "cancelled"
	StepPending    = "pending"
	StepRunning    = "running"
	StepSucceeded  = "succeeded"
	StepFailed     = "failed"
	StepSkipped    = "skipped"
)

const (
	SourceUser      = "user"
	SourceWorkflow  = "workflow"
	SourceSystem    = "system"
	SourceImport    = "import"
	SourceManual    = "manual"
	SourceScheduled = "scheduled"
)

// Workflow is the org-scoped automation header.
type Workflow struct {
	ID                 string
	OrganizationID     string
	ModuleID           *string
	Name               string
	Description        string
	Status             string
	OnActionError      string
	Priority           int
	PublishedVersionID *string
	CreatedBy          *string
	UpdatedBy          *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// Version is an immutable (once published) definition snapshot.
type Version struct {
	ID                 string
	WorkflowID         string
	OrganizationID     string
	Version            int
	State              string
	DefinitionSnapshot json.RawMessage
	Changelog          string
	PublishedAt        *time.Time
	PublishedBy        *string
	CreatedAt          time.Time
}

// Trigger fires a workflow for a version.
type Trigger struct {
	ID             string
	VersionID      string
	OrganizationID string
	Type           string
	Config         json.RawMessage
	CreatedAt      time.Time
}

// Condition is a group or predicate node in the condition tree.
type Condition struct {
	ID             string
	VersionID      string
	OrganizationID string
	ParentID       *string
	NodeType       string
	Logic          *string
	FieldAPIName   *string
	Operator       *string
	Value          json.RawMessage
	SortOrder      int
	CreatedAt      time.Time
}

// Action is an ordered step on a version.
type Action struct {
	ID              string
	VersionID       string
	OrganizationID  string
	SortOrder       int
	Type            string
	Config          json.RawMessage
	MaxRetries      int
	ContinueOnError *bool
	CreatedAt       time.Time
}

// Execution is an append-only run header.
type Execution struct {
	ID             string
	OrganizationID string
	WorkflowID     string
	VersionID      *string
	ModuleID       *string
	RecordID       *string
	TriggerType    string
	Status         string
	Source         string
	Depth          int
	ErrorSummary   *string
	StartedAt      *time.Time
	FinishedAt     *time.Time
	DurationMs     *int
	CreatedAt      time.Time
}

// ExecutionStep logs one action within a run.
type ExecutionStep struct {
	ID             string
	ExecutionID    string
	OrganizationID string
	ActionID       *string
	SortOrder      int
	ActionType     string
	Status         string
	Input          json.RawMessage
	Output         json.RawMessage
	Error          *string
	StartedAt      *time.Time
	FinishedAt     *time.Time
	CreatedAt      time.Time
}

// Template is a built-in or org-cloned starter definition.
type Template struct {
	ID            string
	OrganizationID *string
	Key           string
	Name          string
	Description   string
	ModuleAPIName *string
	Definition    json.RawMessage
	IsBuiltin     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// MatchCandidate is an active workflow with its published version definition.
type MatchCandidate struct {
	Workflow Workflow
	Version  Version
	Triggers []Trigger
	Conditions []Condition
	Actions  []Action
}
