package entity

import "time"

// DeliveryEvent is an append-only audit entry for a notification lifecycle step.
type DeliveryEvent struct {
	ID             string
	OrganizationID string
	NotificationID string
	Event          string
	Provider       *string
	Payload        []byte // JSONB
	CreatedAt      time.Time
}
