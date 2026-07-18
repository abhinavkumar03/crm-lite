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
