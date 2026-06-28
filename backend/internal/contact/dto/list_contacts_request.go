package dto

type ListContactsRequest struct {
	Page   int
	Limit  int
	Search string
}
