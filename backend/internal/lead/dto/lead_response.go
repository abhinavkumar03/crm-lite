package dto

import "time"

type LeadResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Email string `json:"email"`

	Phone string `json:"phone"`

	Company string `json:"company"`

	Status string `json:"status"`

	Notes string `json:"notes"`

	OwnerID string `json:"owner_id"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
