package entity

import "time"

// Role is an organization-scoped named set of permissions and ACL rules.
type Role struct {
	ID             string
	OrganizationID string
	Name           string
	Slug           string
	Description    *string
	IsSystem       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Permission is a global catalog entry.
type Permission struct {
	ID          string
	Key         string
	Category    string
	Description *string
}
