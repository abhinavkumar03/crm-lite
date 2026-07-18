package dto

import (
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
)

// SettingsResponse is the API representation of an organization's settings.
type SettingsResponse struct {
	ID         string                    `json:"id"`
	Name       string                    `json:"name"`
	Slug       string                    `json:"slug"`
	Plan       string                    `json:"plan"`
	General    entity.GeneralSettings    `json:"general"`
	Automation entity.AutomationSettings `json:"automation"`
	UpdatedAt  time.Time                 `json:"updated_at"`
}

// UpdateSettingsRequest is a partial update. Only non-nil sections are applied,
// so the client can save one tab at a time without clobbering the others.
type UpdateSettingsRequest struct {
	Name       *string                    `json:"name" validate:"omitempty,min=1,max=200"`
	General    *entity.GeneralSettings    `json:"general"`
	Automation *entity.AutomationSettings `json:"automation"`
}
