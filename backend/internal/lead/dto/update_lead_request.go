package dto

type UpdateLeadRequest struct {
	Name string `json:"name"`

	Email string `json:"email"`

	Phone string `json:"phone"`

	Company string `json:"company"`

	Status string `json:"status"`

	Notes string `json:"notes"`
}
