package dto

import "time"

// Storage kinds describe where a field's value physically lives.
const (
	StorageColumn = "column" // a real column on a native module's table
	StorageJSONB  = "jsonb"  // a key inside records.data for dynamic modules
)

// FieldOption is a single choice for dropdown/multiselect/radio fields.
type FieldOption struct {
	Label string `json:"label" validate:"required,max=120"`
	Value string `json:"value" validate:"required,max=120"`
}

// StorageDescriptor tells the client (and the runtime record engine) how a
// field is persisted — the metadata-driven persistence strategy.
type StorageDescriptor struct {
	Kind string `json:"kind"` // column | jsonb
	Path string `json:"path"` // column name, or "data.<api_name>"
}

// FieldResponse is the full rendering + persistence metadata for a field.
type FieldResponse struct {
	ID                string            `json:"id"`
	ModuleID          string            `json:"module_id"`
	APIName           string            `json:"api_name"`
	Label             string            `json:"label"`
	FieldType         string            `json:"field_type"`
	IsRequired        bool              `json:"is_required"`
	IsUnique          bool              `json:"is_unique"`
	IsReadOnly        bool              `json:"is_read_only"`
	DefaultValue      *string           `json:"default_value"`
	Placeholder       *string           `json:"placeholder"`
	Description       *string           `json:"description"`
	HelpText          *string           `json:"help_text"`
	MinLength         *int              `json:"min_length"`
	MaxLength         *int              `json:"max_length"`
	Regex             *string           `json:"regex"`
	ValidationMessage *string           `json:"validation_message"`
	Options           []FieldOption     `json:"options"`
	LookupModuleID    *string           `json:"lookup_module_id"`
	SortOrder         int               `json:"sort_order"`
	IsVisible         bool              `json:"is_visible"`
	IsSearchable      bool              `json:"is_searchable"`
	IsFilterable      bool              `json:"is_filterable"`
	IsNullable        bool              `json:"is_nullable"`
	IsIndexed         bool              `json:"is_indexed"`
	IsSystem          bool              `json:"is_system"`
	LockMode          string            `json:"lock_mode"`
	EditableBy        string            `json:"editable_by"`
	ViewableBy        string            `json:"viewable_by"`
	Storage           StorageDescriptor `json:"storage"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// CreateFieldRequest is the payload for adding a field to a module.
type CreateFieldRequest struct {
	APIName           string        `json:"api_name" validate:"required,max=80"`
	Label             string        `json:"label" validate:"required,max=120"`
	FieldType         string        `json:"field_type" validate:"required"`
	IsRequired        bool          `json:"is_required"`
	IsUnique          bool          `json:"is_unique"`
	IsReadOnly        bool          `json:"is_read_only"`
	DefaultValue      *string       `json:"default_value"`
	Placeholder       *string       `json:"placeholder" validate:"omitempty,max=200"`
	Description       *string       `json:"description"`
	HelpText          *string       `json:"help_text"`
	MinLength         *int          `json:"min_length" validate:"omitempty,min=0"`
	MaxLength         *int          `json:"max_length" validate:"omitempty,min=1"`
	Regex             *string       `json:"regex"`
	ValidationMessage *string       `json:"validation_message"`
	Options           []FieldOption `json:"options" validate:"omitempty,dive"`
	LookupModuleID    *string       `json:"lookup_module_id" validate:"omitempty,uuid"`
	IsVisible         *bool         `json:"is_visible"`
	IsSearchable      bool          `json:"is_searchable"`
	IsFilterable      bool          `json:"is_filterable"`
	LockMode          string        `json:"lock_mode" validate:"omitempty,oneof=never after_create always"`
	EditableBy        string        `json:"editable_by" validate:"omitempty,max=40"`
	ViewableBy        string        `json:"viewable_by" validate:"omitempty,max=40"`
	SectionKey        string        `json:"section_key" validate:"omitempty,max=80"`
}

// UpdateFieldRequest is a partial update. api_name and field_type are immutable
// (changing them would orphan stored data), so they are intentionally absent.
type UpdateFieldRequest struct {
	Label             *string       `json:"label" validate:"omitempty,max=120"`
	IsRequired        *bool         `json:"is_required"`
	IsUnique          *bool         `json:"is_unique"`
	IsReadOnly        *bool         `json:"is_read_only"`
	DefaultValue      *string       `json:"default_value"`
	Placeholder       *string       `json:"placeholder" validate:"omitempty,max=200"`
	Description       *string       `json:"description"`
	HelpText          *string       `json:"help_text"`
	MinLength         *int          `json:"min_length" validate:"omitempty,min=0"`
	MaxLength         *int          `json:"max_length" validate:"omitempty,min=1"`
	Regex             *string       `json:"regex"`
	ValidationMessage *string       `json:"validation_message"`
	Options           []FieldOption `json:"options" validate:"omitempty,dive"`
	IsVisible         *bool         `json:"is_visible"`
	IsSearchable      *bool         `json:"is_searchable"`
	IsFilterable      *bool         `json:"is_filterable"`
	LockMode          *string       `json:"lock_mode" validate:"omitempty,oneof=never after_create always"`
	EditableBy        *string       `json:"editable_by" validate:"omitempty,max=40"`
	ViewableBy        *string       `json:"viewable_by" validate:"omitempty,max=40"`
}

// ReorderRequest carries the new ordering for a set of fields.
type ReorderRequest struct {
	Items []ReorderItem `json:"items" validate:"required,min=1,dive"`
}

type ReorderItem struct {
	ID        string `json:"id" validate:"required,uuid"`
	SortOrder int    `json:"sort_order"`
}
