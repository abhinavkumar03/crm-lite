package dto

import "time"

// AnalyzeResult powers the mapping UI: the detected columns, a small sample of
// rows for preview, an auto-suggested column->field mapping, and the total row
// count so the client can warn about large files before committing.
type AnalyzeResult struct {
	Headers          []string            `json:"headers"`
	SampleRows       []map[string]string `json:"sample_rows"`
	SuggestedMapping map[string]string   `json:"suggested_mapping"` // header -> field api_name
	RowCount         int                 `json:"row_count"`
}

// RowErrorDTO mirrors entity.RowError in the API response.
type RowErrorDTO struct {
	Row     int    `json:"row"`
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// ImportResponse is the API representation of an import job.
type ImportResponse struct {
	ID            string            `json:"id"`
	ModuleID      string            `json:"module_id"`
	Filename      string            `json:"filename"`
	Status        string            `json:"status"`
	Mapping       map[string]string `json:"mapping"`
	TotalRows     int               `json:"total_rows"`
	ProcessedRows int               `json:"processed_rows"`
	SuccessRows   int               `json:"success_rows"`
	ErrorRows     int               `json:"error_rows"`
	Errors        []RowErrorDTO     `json:"errors"`
	CreatedBy     *string           `json:"created_by"`
	StartedAt     *time.Time        `json:"started_at"`
	FinishedAt    *time.Time        `json:"finished_at"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

// ListQuery is the parsed, sanitized set of list parameters.
type ListQuery struct {
	Page     int
	PageSize int
	Status   string
}

// ListResult is a paginated collection of import jobs.
type ListResult struct {
	Imports    []ImportResponse `json:"imports"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	Total      int              `json:"total"`
	TotalPages int              `json:"total_pages"`
}
