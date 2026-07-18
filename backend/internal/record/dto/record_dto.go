package dto

import "time"

// FilterClause is a single dynamic filter applied to a JSONB field.
type FilterClause struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    any    `json:"value"`
}

// Supported filter operators (evaluated server-side against the JSONB payload).
const (
	OpEquals      = "eq"
	OpNotEquals   = "ne"
	OpContains    = "contains"
	OpGreaterThan = "gt"
	OpLessThan    = "lt"
	OpGreaterEq   = "gte"
	OpLessEq      = "lte"
	OpIn          = "in"
)

// ListQuery is the parsed, sanitized set of list parameters.
type ListQuery struct {
	Page     int
	PageSize int
	Search   string
	Sort     string
	Order    string
	Filters  []FilterClause
	Expand   bool
	// SkipTotal omits the COUNT(*) query. Used by export paging where only
	// "is there another page?" matters (short page = done).
	SkipTotal bool
}

// RelationRef is a resolved lookup/user reference (id + human label).
type RelationRef struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// RecordResponse is the API representation of a record.
type RecordResponse struct {
	ID        string                 `json:"id"`
	ModuleID  string                 `json:"module_id"`
	Data      map[string]any         `json:"data"`
	OwnerID   *string                `json:"owner_id"`
	CreatedBy *string                `json:"created_by"`
	UpdatedBy *string                `json:"updated_by"`
	Relations map[string]RelationRef `json:"relations,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// ListResult is a paginated collection of records.
type ListResult struct {
	Records    []RecordResponse `json:"records"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	Total      int              `json:"total"`
	TotalPages int              `json:"total_pages"`
}

// CreateRecordRequest creates a record. Data is validated against the module's
// field metadata + validation rules before it is persisted.
type CreateRecordRequest struct {
	Data    map[string]any `json:"data" validate:"required"`
	OwnerID *string        `json:"owner_id"`
}

// UpdateRecordRequest replaces a record's data payload.
type UpdateRecordRequest struct {
	Data    map[string]any `json:"data" validate:"required"`
	OwnerID *string        `json:"owner_id"`
}
