package dto

import "time"

type CreateCallLogRequest struct {
	Direction string `json:"direction" binding:"required"`

	Status string `json:"status" binding:"required"`

	DurationSeconds int `json:"duration_seconds"`

	Summary string `json:"summary"`

	FollowUpAt *time.Time `json:"follow_up_at"`
}
