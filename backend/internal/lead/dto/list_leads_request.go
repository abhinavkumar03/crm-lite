package dto

type ListLeadRequest struct {
	Page      int
	Limit     int
	Search    string
	Status    string
	SortBy    string
	SortOrder string
}
