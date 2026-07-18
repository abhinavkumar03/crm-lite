package entity

import "time"

// Notification lifecycle states.
const (
	StatusQueued = "queued"
	StatusSent   = "sent"
	StatusFailed = "failed"
)

// Notification is a durable outbound-message record. Content is rendered at
// creation time; the worker only transitions status and stamps the provider.
type Notification struct {
	ID             string
	OrganizationID string
	Channel        string
	Recipient      string
	Subject        *string
	Body           string
	Template       *string
	Data           []byte // JSONB
	Status         string
	Provider       *string
	Error          *string
	EntityType     *string
	EntityID       *string
	CreatedBy      *string
	SentAt         *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
