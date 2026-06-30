package entity

import "time"

type EntityType string

const (
	EntityLead EntityType = "LEAD"

	EntityContact EntityType = "CONTACT"

	EntityTask EntityType = "TASK"
)

const (
	ActionLeadCreated = "LEAD_CREATED"

	ActionLeadUpdated = "LEAD_UPDATED"

	ActionLeadStatusChanged = "LEAD_STATUS_CHANGED"

	ActionLeadDeleted = "LEAD_DELETED"

	ActionContactCreated = "CONTACT_ADDED"

	ActionContactUpdated = "CONTACT_UPDATED"

	ActionContactDeleted = "CONTACT_DELETED"

	ActionNoteAdded = "NOTE_ADDED"

	ActionNoteUpdated = "NOTE_UPDATED"

	ActionNoteDeleted = "NOTE_DELETED"

	ActionCallLogged = "CALL_LOGGED"

	ActionCallUpdated = "CALL_UPDATED"

	ActionCallDeleted = "CALL_DELETED"

	ActionAttachmentAdded = "ATTACHMENT_ADDED"

	ActionAttachmentDeleted = "ATTACHMENT_DELETED"

	ActionTaskCreated = "TASK_CREATED"

	ActionTaskUpdated = "TASK_UPDATED"

	ActionTaskDeleted = "TASK_DELETED"

	ActionTaskCompleted = "TASK_COMPLETED"
)

type Activity struct {
	ID string

	EntityType EntityType

	EntityID string

	Action string

	Description string

	PerformedBy string

	Metadata []byte

	CreatedAt time.Time
}
