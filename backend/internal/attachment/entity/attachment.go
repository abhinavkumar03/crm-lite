package entity

import "time"

type EntityType string

const (
	EntityLead    EntityType = "LEAD"
	EntityContact EntityType = "CONTACT"
	EntityTask    EntityType = "TASK"
)

type Attachment struct {
	ID string

	EntityType EntityType

	EntityID string

	FileName string

	FileURL string

	PublicID string

	ResourceType string

	FileSize int64

	UploadedBy string

	CreatedAt time.Time
}
