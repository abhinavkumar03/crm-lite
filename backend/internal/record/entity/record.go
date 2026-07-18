package entity

import "time"

// Record is a single row of a dynamic module, stored generically in the records
// table. All user-defined field values live in the Data JSONB column keyed by
// each field's api_name.
type Record struct {
	ID             string
	OrganizationID string
	ModuleID       string
	Data           []byte // JSONB
	OwnerID        *string
	CreatedBy      *string
	UpdatedBy      *string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
