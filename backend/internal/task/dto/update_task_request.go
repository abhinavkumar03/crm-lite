package dto

type UpdateTaskRequest struct {
	Title string `json:"title"`

	Description string `json:"description"`

	Status string `json:"status"`

	LeadID *string `json:"lead_id"`

	ContactID *string `json:"contact_id"`

	DueDate *string `json:"due_date"`
}
