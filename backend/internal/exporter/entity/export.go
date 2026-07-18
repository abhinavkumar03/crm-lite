package entity

import "time"

// Export job lifecycle states.
const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

// Output formats.
const (
	FormatCSV  = "csv"
	FormatXLSX = "xlsx"
)

// ExportJob is the durable record of an export. For a completed job the generated
// file is stored inline in Content so it can be re-downloaded from history.
type ExportJob struct {
	ID             string
	OrganizationID string
	ModuleID       string
	Filename       string
	Format         string
	Status         string
	Columns        []byte // JSONB: ["api_name", ...]
	Filters        []byte // JSONB: [{field, operator, value}, ...]
	Options        []byte // JSONB
	RowCount       int
	ByteSize       int
	Content        []byte // BYTEA (nil until completed)
	Error          *string
	CreatedBy      *string
	StartedAt      *time.Time
	FinishedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ExportTemplate is a reusable export configuration for a module.
type ExportTemplate struct {
	ID             string
	OrganizationID string
	ModuleID       string
	Name           string
	Format         string
	Columns        []byte // JSONB
	Filters        []byte // JSONB
	Sort           []byte // JSONB
	CreatedBy      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
