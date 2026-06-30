package dto

import "time"

type ActivityResponse struct {
	ID string `json:"id"`

	Action string `json:"action"`

	Description string `json:"description"`

	PerformedBy string `json:"performed_by"`

	Metadata any `json:"metadata,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}
