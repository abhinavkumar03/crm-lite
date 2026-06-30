package dto

type CreateNoteRequest struct {
	EntityType string `json:"entity_type" binding:"required"`
	EntityID   string `json:"entity_id" binding:"required"`
	Note       string `json:"note" binding:"required,max=5000"`
}
