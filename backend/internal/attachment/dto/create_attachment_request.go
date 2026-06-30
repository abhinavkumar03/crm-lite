package dto

import "time"

type CreateAttachmentRequest struct {
	FileName string `json:"file_name" binding:"required"`

	FileURL string `json:"file_url" binding:"required,url"`

	PublicID string `json:"public_id" binding:"required"`

	ResourceType string `json:"resource_type" binding:"required"`

	FileSize int64 `json:"file_size"`
}

type AttachmentResponse struct {
	ID string `json:"id"`

	EntityType string `json:"entity_type"`

	EntityID string `json:"entity_id"`

	FileName string `json:"file_name"`

	FileURL string `json:"file_url"`

	PublicID string `json:"public_id"`

	ResourceType string `json:"resource_type"`

	FileSize int64 `json:"file_size"`

	UploadedBy string `json:"uploaded_by"`

	CreatedAt time.Time `json:"created_at"`
}
