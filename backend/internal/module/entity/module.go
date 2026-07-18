package entity

import "time"

// Storage strategies for a module.
const (
	StorageNative  = "native"
	StorageDynamic = "dynamic"
)

// Module is a metadata-defined object type.
type Module struct {
	ID               string
	OrganizationID   string
	APIName          string
	SingularLabel    string
	PluralLabel      string
	Description      *string
	Icon             *string
	Color            *string
	StorageStrategy  string
	NativeTable      *string
	IsSystem         bool
	IsEnabled        bool
	IsVisibleSidebar bool
	SortOrder        int
	DefaultSortField string
	DefaultSortOrder string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// SortPosition is a single (id, sort_order) pair used for reordering.
type SortPosition struct {
	ID        string
	SortOrder int
}
