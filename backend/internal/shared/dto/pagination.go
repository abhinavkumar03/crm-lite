package dto

type PaginationRequest struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}
