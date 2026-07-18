package dto

import "time"

// UpdateProgressRequest is the single write path used by the client to advance,
// complete or skip a tour. All fields are optional so a caller can update just
// the cursor (current_step) or just the lifecycle (status) without clobbering
// the rest.
type UpdateProgressRequest struct {
	TourKey        string   `json:"tour_key" validate:"omitempty,max=80"`
	Status         string   `json:"status" validate:"omitempty,oneof=active completed skipped"`
	CurrentStep    *int     `json:"current_step" validate:"omitempty,min=0"`
	CompletedSteps []string `json:"completed_steps" validate:"omitempty,dive,max=120"`
}

// RestartRequest resets a tour back to the beginning.
type RestartRequest struct {
	TourKey string `json:"tour_key" validate:"omitempty,max=80"`
}

// ProgressResponse is the API representation of a user's tour progress.
type ProgressResponse struct {
	TourKey        string     `json:"tour_key"`
	Status         string     `json:"status"`
	CurrentStep    int        `json:"current_step"`
	CompletedSteps []string   `json:"completed_steps"`
	StartedAt      time.Time  `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
