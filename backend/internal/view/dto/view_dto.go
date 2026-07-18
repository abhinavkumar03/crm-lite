package dto

import "time"

// ViewFilter is a single filter clause stored on a view.
type ViewFilter struct {
	Field    string `json:"field" validate:"required,max=80"`
	Operator string `json:"operator" validate:"required,max=20"`
	Value    any    `json:"value"`
}

// ViewSort is the sort configuration for a view.
type ViewSort struct {
	Field string `json:"field" validate:"omitempty,max=80"`
	Order string `json:"order" validate:"omitempty,oneof=asc desc"`
}

// ViewResponse is the API representation of a saved view.
type ViewResponse struct {
	ID        string       `json:"id"`
	ModuleID  string       `json:"module_id"`
	Name      string       `json:"name"`
	Columns   []string     `json:"columns"`
	Filters   []ViewFilter `json:"filters"`
	Sort      ViewSort     `json:"sort"`
	IsDefault bool         `json:"is_default"`
	IsPublic  bool         `json:"is_public"`
	OwnerID   *string      `json:"owner_id"`
	IsOwner   bool         `json:"is_owner"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// CreateViewRequest is the payload for saving a new view.
type CreateViewRequest struct {
	Name     string       `json:"name" validate:"required,max=120"`
	Columns  []string     `json:"columns" validate:"required,min=1,dive,required,max=80"`
	Filters  []ViewFilter `json:"filters" validate:"omitempty,dive"`
	Sort     *ViewSort    `json:"sort" validate:"omitempty"`
	IsPublic *bool        `json:"is_public"`
}

// UpdateViewRequest is a partial update of a saved view.
type UpdateViewRequest struct {
	Name     *string      `json:"name" validate:"omitempty,max=120"`
	Columns  []string     `json:"columns" validate:"omitempty,min=1,dive,required,max=80"`
	Filters  []ViewFilter `json:"filters" validate:"omitempty,dive"`
	Sort     *ViewSort    `json:"sort" validate:"omitempty"`
	IsPublic *bool        `json:"is_public"`
}
