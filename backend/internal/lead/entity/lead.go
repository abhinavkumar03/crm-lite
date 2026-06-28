package entity

import "time"

type Lead struct {
	ID string `db:"id" json:"id"`

	OwnerID string `db:"owner_id" json:"owner_id"`

	Name string `db:"name" json:"name"`

	Email string `db:"email" json:"email"`

	Phone string `db:"phone" json:"phone"`

	Company string `db:"company" json:"company"`

	Status string `db:"status" json:"status"`

	Notes string `db:"notes" json:"notes"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
