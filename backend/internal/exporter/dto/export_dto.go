package dto

import (
	"time"

	recorddto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
)

// ExportSpec is the configuration shared by sync downloads, async exports and
// saved templates: what to export (columns/filters/search/sort) and how (format).
// An empty Columns list means "all visible fields". Filters reuse the record
// runtime's clause shape so the same query engine (Phase 10) drives exports.
type ExportSpec struct {
	Format  string                   `json:"format"`
	Columns []string                 `json:"columns"`
	Filters []recorddto.FilterClause `json:"filters"`
	Search  string                   `json:"search"`
	Sort    string                   `json:"sort"`
	Order   string                   `json:"order"`
	Expand  bool                     `json:"expand"`
}

// ExportResponse is the API representation of an export job (never includes the
// file bytes; those are served by the dedicated download endpoint).
type ExportResponse struct {
	ID         string     `json:"id"`
	ModuleID   string     `json:"module_id"`
	Filename   string     `json:"filename"`
	Format     string     `json:"format"`
	Status     string     `json:"status"`
	Columns    []string   `json:"columns"`
	RowCount   int        `json:"row_count"`
	ByteSize   int        `json:"byte_size"`
	Error      *string    `json:"error"`
	CreatedBy  *string    `json:"created_by"`
	StartedAt  *time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ListQuery is the parsed, sanitized set of list parameters.
type ListQuery struct {
	Page     int
	PageSize int
	Status   string
}

// ListResult is a paginated collection of export jobs.
type ListResult struct {
	Exports    []ExportResponse `json:"exports"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	Total      int              `json:"total"`
	TotalPages int              `json:"total_pages"`
}

// --- Templates -------------------------------------------------------------

// TemplateSort is the persisted sort of a template.
type TemplateSort struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

// CreateTemplateRequest saves a reusable export configuration.
type CreateTemplateRequest struct {
	Name    string                   `json:"name" validate:"required,max=120"`
	Format  string                   `json:"format" validate:"omitempty,oneof=csv xlsx"`
	Columns []string                 `json:"columns"`
	Filters []recorddto.FilterClause `json:"filters"`
	Sort    *TemplateSort            `json:"sort"`
}

// UpdateTemplateRequest patches a template; nil fields are left unchanged.
type UpdateTemplateRequest struct {
	Name    *string                  `json:"name" validate:"omitempty,max=120"`
	Format  *string                  `json:"format" validate:"omitempty,oneof=csv xlsx"`
	Columns []string                 `json:"columns"`
	Filters []recorddto.FilterClause `json:"filters"`
	Sort    *TemplateSort            `json:"sort"`
}

// TemplateResponse is the API representation of an export template.
type TemplateResponse struct {
	ID        string                   `json:"id"`
	ModuleID  string                   `json:"module_id"`
	Name      string                   `json:"name"`
	Format    string                   `json:"format"`
	Columns   []string                 `json:"columns"`
	Filters   []recorddto.FilterClause `json:"filters"`
	Sort      TemplateSort             `json:"sort"`
	CreatedBy *string                  `json:"created_by"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}
