package dto

import "time"

type UpdateCallLogRequest struct {
	Direction string `json:"direction"`

	Status string `json:"status"`

	DurationSeconds int `json:"duration_seconds"`

	Summary string `json:"summary"`

	FollowUpAt *time.Time `json:"follow_up_at"`
}
