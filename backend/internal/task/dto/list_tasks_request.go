package dto

type ListTasksRequest struct {
	Page      int
	Limit     int
	Search    string
	Status    string
	SortBy    string
	SortOrder string
}
