package dto

type UpdateContactRequest struct {
	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	Email string `json:"email"`

	Phone string `json:"phone"`

	Company string `json:"company"`

	JobTitle string `json:"job_title"`

	Notes string `json:"notes"`
}
