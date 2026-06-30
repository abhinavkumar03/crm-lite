package dto

import "time"

type TaskResponse struct {
	ID string `json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	Status string `json:"status"`

	DueDate *string `json:"due_date,omitempty"`

	LeadID *string `json:"lead_id,omitempty"`

	ContactID *string `json:"contact_id,omitempty"`

	OwnerID string `json:"owner_id"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
