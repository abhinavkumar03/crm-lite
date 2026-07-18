package entity

import "time"

// Tour lifecycle states.
const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusSkipped   = "skipped"
)

// DefaultTourKey is the key of the application-wide onboarding tour. The schema
// supports multiple named tours per user; the client currently drives this one.
const DefaultTourKey = "app"

// TourProgress is a user's durable progress through a named guided tour. The
// step catalogue lives on the client, so this record only tracks position and
// lifecycle. CompletedSteps holds the client step keys the user has seen.
type TourProgress struct {
	ID             string
	OrganizationID string
	UserID         string
	TourKey        string
	Status         string
	CurrentStep    int
	CompletedSteps []string
	StartedAt      time.Time
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
