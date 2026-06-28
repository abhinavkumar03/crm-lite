package dto

type ListLeadsRequest struct {
	Page int

	Limit int

	Search string

	Status string
}
