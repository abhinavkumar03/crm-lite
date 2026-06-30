package dto

import "time"

type NoteUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NoteResponse struct {
	ID         string    `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   string    `json:"entity_id"`
	Note       string    `json:"note"`
	CreatedBy  string    `json:"created_by"`
	UpdatedBy  *string   `json:"updated_by,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	User NoteUser `json:"user"`
}
