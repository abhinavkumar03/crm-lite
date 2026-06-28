package entity

import "time"

type Contact struct {
	ID string `db:"id" json:"id"`

	OwnerID string `db:"owner_id" json:"owner_id"`

	FirstName string `db:"first_name" json:"first_name"`

	LastName string `db:"last_name" json:"last_name"`

	Email string `db:"email" json:"email"`

	Phone string `db:"phone" json:"phone"`

	Company string `db:"company" json:"company"`

	JobTitle string `db:"job_title" json:"job_title"`

	Notes string `db:"notes" json:"notes"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`

	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
