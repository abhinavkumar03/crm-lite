package dto

import (
	"encoding/json"
	"time"
)

type LayoutResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Type      string          `json:"layout_type"`
	IsDefault bool            `json:"is_default"`
	Config    json.RawMessage `json:"config"`
}

// LayoutSection is one group of fields on a detail (or form) layout.
type LayoutSection struct {
	Key         string   `json:"key" validate:"required,max=80"`
	Label       string   `json:"label" validate:"required,max=120"`
	Description string   `json:"description,omitempty"`
	Order       int      `json:"order,omitempty"`
	Collapsed   bool     `json:"collapsed,omitempty"`
	Columns     int      `json:"columns,omitempty"` // 1 | 2 | 3
	Fields      []string `json:"fields"`
}

// UpdateDetailLayoutRequest replaces the default detail layout section config.
type UpdateDetailLayoutRequest struct {
	Sections []LayoutSection `json:"sections" validate:"required,min=1,dive"`
	Tabs     []string        `json:"tabs"`
}

// --- Form layout (hydrated) -------------------------------------------------

// FormFieldOption is a choice for dropdown/multiselect fields in form metadata.
type FormFieldOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// FormFieldValidationRules is a compact validation object for form rendering.
type FormFieldValidationRules struct {
	Min   *int    `json:"min,omitempty"`
	Max   *int    `json:"max,omitempty"`
	Regex *string `json:"regex,omitempty"`
}

// FormLayoutField is a hydrated field inside a form section.
type FormLayoutField struct {
	ID              string                   `json:"id"`
	Key             string                   `json:"key"`
	Label           string                   `json:"label"`
	Type            string                   `json:"type"`
	Required        bool                     `json:"required"`
	Editable        bool                     `json:"editable"`
	Locked          bool                     `json:"locked"`
	DisplayOrder    int                      `json:"display_order"`
	Placeholder     *string                  `json:"placeholder,omitempty"`
	Description     *string                  `json:"description,omitempty"`
	DefaultValue    *string                  `json:"default_value,omitempty"`
	ValidationRules FormFieldValidationRules `json:"validation_rules"`
	Options         []FormFieldOption        `json:"options"`
	LookupModuleID  *string                  `json:"lookup_module_id,omitempty"`
	LockMode        string                   `json:"lock_mode"`
}

// FormLayoutSection is a hydrated section for create/edit forms.
type FormLayoutSection struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description,omitempty"`
	Order       int               `json:"order"`
	Collapsed   bool              `json:"collapsed"`
	Columns     int               `json:"columns"`
	Fields      []FormLayoutField `json:"fields"`
}

// FormLayoutResponse is returned by GET /layouts/form.
type FormLayoutResponse struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	LayoutType string              `json:"layout_type"`
	IsDefault  bool                `json:"is_default"`
	Mode       string              `json:"mode"`
	Sections   []FormLayoutSection `json:"sections"`
}

// UpdateFormLayoutRequest replaces the default form layout section config.
type UpdateFormLayoutRequest struct {
	Sections []LayoutSection `json:"sections" validate:"required,min=1,dive"`
}

// FormReorderRequest reorders fields within one section.
type FormReorderRequest struct {
	SectionID string            `json:"section_id" validate:"required,max=80"`
	Fields    []FormReorderItem `json:"fields" validate:"required,min=1,dive"`
}

type FormReorderItem struct {
	FieldID string `json:"field_id" validate:"required"` // field UUID or api_name
	Order   int    `json:"order"`
}

// CreateSectionRequest adds a section to the form layout.
type CreateSectionRequest struct {
	Key         string `json:"key" validate:"required,max=80"`
	Label       string `json:"label" validate:"required,max=120"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Columns     int    `json:"columns" validate:"omitempty,min=1,max=3"`
	Collapsed   bool   `json:"collapsed"`
}

// UpdateSectionRequest updates section metadata (not field membership).
type UpdateSectionRequest struct {
	Label       *string `json:"label" validate:"omitempty,max=120"`
	Description *string `json:"description" validate:"omitempty,max=500"`
	Columns     *int    `json:"columns" validate:"omitempty,min=1,max=3"`
	Collapsed   *bool   `json:"collapsed"`
	Order       *int    `json:"order"`
}

// --- List layout ------------------------------------------------------------

const ActionsColumnKey = "_actions"

// ListColumn is one column in the org default list layout.
// field_id and label are hydrated on GET (not persisted in layout JSON).
// locked is derived on hydrate; width is optional/future-ready.
type ListColumn struct {
	FieldKey   string `json:"field_key" validate:"required,max=80"`
	FieldID    string `json:"field_id,omitempty"`
	Label      string `json:"label,omitempty"`
	Visible    bool   `json:"visible"`
	Order      int    `json:"order"`
	Sortable   bool   `json:"sortable"`
	Searchable bool   `json:"searchable"`
	System     bool   `json:"system"`
	Locked     bool   `json:"locked"`
	Width      *int   `json:"width,omitempty"`
}

// ListLayoutResponse is returned by GET /layouts/list (visible columns only for runtime).
type ListLayoutResponse struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	LayoutType string       `json:"layout_type"`
	IsDefault  bool         `json:"is_default"`
	Columns    []ListColumn `json:"columns"`
}

// UpdateListLayoutRequest replaces the default list column config.
type UpdateListLayoutRequest struct {
	Columns []ListColumn `json:"columns" validate:"required,min=1,dive"`
}

// ListReorderRequest updates column order only.
type ListReorderRequest struct {
	Columns []ListReorderItem `json:"columns" validate:"required,min=1,dive"`
}

type ListReorderItem struct {
	FieldKey string `json:"field_key" validate:"required,max=80"`
	Order    int    `json:"order"`
}

// ListToggleRequest flips visibility for one column.
type ListToggleRequest struct {
	FieldKey string `json:"field_key" validate:"required,max=80"`
	Visible  *bool  `json:"visible" validate:"required"`
}

type NoteResponse struct {
	ID         string    `json:"id"`
	Title      *string   `json:"title,omitempty"`
	Body       string    `json:"body"`
	CreatedBy  string    `json:"created_by"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title *string `json:"title"`
	Body  string  `json:"body" validate:"required,min=1"`
}

type AttachmentResponse struct {
	ID           string    `json:"id"`
	FileName     string    `json:"file_name"`
	FileURL      string    `json:"file_url"`
	PublicID     string    `json:"public_id"`
	ResourceType *string   `json:"resource_type,omitempty"`
	FileSize     *int64    `json:"file_size,omitempty"`
	UploadedBy   string    `json:"uploaded_by"`
	UploaderName string    `json:"uploader_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateAttachmentRequest struct {
	FileName     string  `json:"file_name" validate:"required,min=1,max=255"`
	FileURL      string  `json:"file_url" validate:"required,url"`
	PublicID     string  `json:"public_id" validate:"required,min=1"`
	ResourceType *string `json:"resource_type"`
	FileSize     *int64  `json:"file_size"`
}

type ActivityResponse struct {
	ID          string          `json:"id"`
	Action      string          `json:"action"`
	Description string          `json:"description"`
	PerformedBy string          `json:"performed_by"`
	ActorName   string          `json:"actor_name"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type RelatedDescriptorResponse struct {
	ChildModuleID    string `json:"child_module_id"`
	ChildModuleName  string `json:"child_module_name"`
	ChildAPIName     string `json:"child_api_name"`
	LookupFieldAPI   string `json:"lookup_field_api_name"`
	LookupFieldLabel string `json:"lookup_field_label"`
}
