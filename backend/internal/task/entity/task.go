package entity

import "time"

type Task struct {
	ID string `db:"id" json:"id"`

	OwnerID string `db:"owner_id" json:"owner_id"`

	LeadID *string `db:"lead_id" json:"lead_id,omitempty"`

	ContactID *string `db:"contact_id" json:"contact_id,omitempty"`

	Title string `db:"title" json:"title"`

	Description string `db:"description" json:"description"`

	Status string `db:"status" json:"status"`

	DueDate *time.Time `db:"due_date" json:"due_date,omitempty"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
