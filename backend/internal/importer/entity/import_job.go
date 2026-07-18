package entity

import "time"

// Import job lifecycle states.
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// RowError captures a single failure for one source row (1-based). Field is set
// for validation failures and empty for persistence-level errors.
type RowError struct {
	Row     int    `json:"row"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ImportJob is the durable record of a bulk import. SourceRows stages the parsed
// file so the worker can process, retry and report entirely from the database.
type ImportJob struct {
	ID             string
	OrganizationID string
	ModuleID       string
	Filename       string
	Status         string
	Mapping        []byte // JSONB: source header -> field api_name
	Options        []byte // JSONB
	SourceRows     []byte // JSONB: []map[string]string
	TotalRows      int
	ProcessedRows  int
	SuccessRows    int
	ErrorRows      int
	Errors         []byte // JSONB: []RowError
	CreatedBy      *string
	StartedAt      *time.Time
	FinishedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
