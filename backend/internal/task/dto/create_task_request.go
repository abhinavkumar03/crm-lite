package dto

type CreateTaskRequest struct {
	Title string `json:"title" validate:"required,max=200"`

	Description string `json:"description"`

	LeadID *string `json:"lead_id"`

	ContactID *string `json:"contact_id"`

	DueDate *string `json:"due_date"`
}
