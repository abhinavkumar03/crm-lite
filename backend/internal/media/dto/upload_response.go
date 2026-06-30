package dto

type UploadResponse struct {
	URL string `json:"url"`

	PublicID string `json:"public_id"`

	ResourceType string `json:"resource_type"`

	Bytes int `json:"bytes"`

	Format string `json:"format"`
}
