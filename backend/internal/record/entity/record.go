package entity

import "time"

// Record is a single row of a dynamic module.
type Record struct {
	ID             string
	OrganizationID string
	ModuleID       string
	Data           []byte
	OwnerID        *string
	AssignedTo     *string
	TeamID         *string
	DepartmentID   *string
	Visibility     string
	CreatedBy      *string
	UpdatedBy      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
