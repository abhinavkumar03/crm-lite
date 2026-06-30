package entity

import "time"

type EntityType string

const (
	EntityLead    EntityType = "LEAD"
	EntityContact EntityType = "CONTACT"
	EntityTask    EntityType = "TASK"
)

type CallDirection string

const (
	Incoming CallDirection = "INCOMING"
	Outgoing CallDirection = "OUTGOING"
)

type CallStatus string

const (
	Completed CallStatus = "COMPLETED"
	Missed    CallStatus = "MISSED"
	NoAnswer  CallStatus = "NO_ANSWER"
	Busy      CallStatus = "BUSY"
	Voicemail CallStatus = "VOICEMAIL"
	Cancelled CallStatus = "CANCELLED"
)

type CallLog struct {
	ID string

	EntityType EntityType

	EntityID string

	Direction CallDirection

	Status CallStatus

	DurationSeconds int

	Summary string

	FollowUpAt *time.Time

	CreatedBy string

	UpdatedBy *string

	CreatedAt time.Time

	UpdatedAt time.Time
}
