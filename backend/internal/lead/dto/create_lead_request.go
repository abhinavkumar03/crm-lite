package dto

type CreateLeadRequest struct {
	Name string `json:"name" validate:"required,max=150"`

	Email string `json:"email" validate:"required,email"`

	Phone string `json:"phone"`

	Company string `json:"company"`

	Notes string `json:"notes"`
}
