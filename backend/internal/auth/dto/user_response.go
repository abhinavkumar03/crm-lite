package dto

type UserResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Email string `json:"email"`

	OrganizationID   string `json:"organization_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
}
