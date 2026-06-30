package entity

import "time"

type EntityType string

const (
	EntityLead    EntityType = "LEAD"
	EntityContact EntityType = "CONTACT"
	EntityTask    EntityType = "TASK"
)

type Note struct {
	ID         string
	EntityType EntityType
	EntityID   string
	Note       string
	CreatedBy  string
	UpdatedBy  *string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
