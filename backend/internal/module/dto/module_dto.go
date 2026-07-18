package dto

import "time"

// ModuleResponse is the API representation of a module.
type ModuleResponse struct {
	ID               string    `json:"id"`
	APIName          string    `json:"api_name"`
	SingularLabel    string    `json:"singular_label"`
	PluralLabel      string    `json:"plural_label"`
	Description      *string   `json:"description"`
	Icon             *string   `json:"icon"`
	Color            *string   `json:"color"`
	StorageStrategy  string    `json:"storage_strategy"`
	NativeTable      *string   `json:"native_table"`
	IsSystem         bool      `json:"is_system"`
	IsEnabled        bool      `json:"is_enabled"`
	IsVisibleSidebar bool      `json:"is_visible_sidebar"`
	SortOrder        int       `json:"sort_order"`
	DefaultSortField string    `json:"default_sort_field"`
	DefaultSortOrder string    `json:"default_sort_order"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// NavigationItem is the trimmed representation used to build the sidebar.
type NavigationItem struct {
	ID            string  `json:"id"`
	APIName       string  `json:"api_name"`
	SingularLabel string  `json:"singular_label"`
	PluralLabel   string  `json:"plural_label"`
	Icon          *string `json:"icon"`
	Color         *string `json:"color"`
	SortOrder     int     `json:"sort_order"`
}

// CreateModuleRequest is the payload for creating a (dynamic) module.
type CreateModuleRequest struct {
	APIName          string  `json:"api_name" validate:"required,max=80"`
	SingularLabel    string  `json:"singular_label" validate:"required,max=80"`
	PluralLabel      string  `json:"plural_label" validate:"required,max=80"`
	Description      *string `json:"description"`
	Icon             *string `json:"icon" validate:"omitempty,max=60"`
	Color            *string `json:"color" validate:"omitempty,max=20"`
	IsVisibleSidebar *bool   `json:"is_visible_sidebar"`
	DefaultSortField *string `json:"default_sort_field" validate:"omitempty,max=80"`
	DefaultSortOrder *string `json:"default_sort_order" validate:"omitempty,oneof=asc desc"`
}

// UpdateModuleRequest is a partial update; nil fields are left unchanged.
type UpdateModuleRequest struct {
	SingularLabel    *string `json:"singular_label" validate:"omitempty,max=80"`
	PluralLabel      *string `json:"plural_label" validate:"omitempty,max=80"`
	Description      *string `json:"description"`
	Icon             *string `json:"icon" validate:"omitempty,max=60"`
	Color            *string `json:"color" validate:"omitempty,max=20"`
	IsVisibleSidebar *bool   `json:"is_visible_sidebar"`
	DefaultSortField *string `json:"default_sort_field" validate:"omitempty,max=80"`
	DefaultSortOrder *string `json:"default_sort_order" validate:"omitempty,oneof=asc desc"`
}

// SetStatusRequest toggles a module's enabled state.
type SetStatusRequest struct {
	Enabled *bool `json:"enabled" validate:"required"`
}

// ReorderRequest carries the new ordering for a set of modules.
type ReorderRequest struct {
	Items []ReorderItem `json:"items" validate:"required,min=1,dive"`
}

type ReorderItem struct {
	ID        string `json:"id" validate:"required,uuid"`
	SortOrder int    `json:"sort_order"`
}
