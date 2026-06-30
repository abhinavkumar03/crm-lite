package dto

import "time"

type CallLogResponse struct {
	ID string `json:"id"`

	EntityType string `json:"entity_type"`

	EntityID string `json:"entity_id"`

	Direction string `json:"direction"`

	Status string `json:"status"`

	DurationSeconds int `json:"duration_seconds"`

	Summary string `json:"summary"`

	FollowUpAt *time.Time `json:"follow_up_at,omitempty"`

	CreatedBy string `json:"created_by"`

	UpdatedBy *string `json:"updated_by,omitempty"`

	CreatedAt time.Time `json:"created_at"`

	UpdatedAt time.Time `json:"updated_at"`
}
