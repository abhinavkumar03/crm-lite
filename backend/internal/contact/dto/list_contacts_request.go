package dto

type ListContactsRequest struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Search    string `json:"search"`
	Company   string `json:"company"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}
