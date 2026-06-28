package dto

type TaskResponse struct {
	ID string `json:"id"`

	Title string `json:"title"`

	Description string `json:"description"`

	Status string `json:"status"`

	LeadID *string `json:"lead_id,omitempty"`

	ContactID *string `json:"contact_id,omitempty"`

	DueDate *string `json:"due_date,omitempty"`
}
