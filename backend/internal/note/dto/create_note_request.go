package dto

type CreateNoteRequest struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Note       string `json:"note" binding:"required,max=5000"`
}
