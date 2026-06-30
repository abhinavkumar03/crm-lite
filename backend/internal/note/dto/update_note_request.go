package dto

type UpdateNoteRequest struct {
	Note string `json:"note" binding:"required,max=5000"`
}
