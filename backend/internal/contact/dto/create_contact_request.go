package dto

type CreateContactRequest struct {
	FirstName string `json:"first_name" validate:"required,max=100"`

	LastName string `json:"last_name"`

	Email string `json:"email" validate:"omitempty,email"`

	Phone string `json:"phone"`

	Company string `json:"company"`

	JobTitle string `json:"job_title"`

	Notes string `json:"notes"`
}
