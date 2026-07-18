package entity

import "time"

// View is a saved table configuration (visible columns, filters, sort) for a
// module. Views can be public (shared across the org) or private to the owner.
type View struct {
	ID             string
	OrganizationID string
	ModuleID       string
	Name           string
	Columns        []byte // JSONB: ["api_name", ...]
	Filters        []byte // JSONB: [{field, operator, value}, ...]
	Sort           []byte // JSONB: {field, order}
	IsDefault      bool
	IsPublic       bool
	OwnerID        *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
