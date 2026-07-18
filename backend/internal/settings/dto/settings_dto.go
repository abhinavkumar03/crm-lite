package dto

import (
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/settings/entity"
)

// SettingsResponse is the API representation of an organization's settings.
type SettingsResponse struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	Slug            string                    `json:"slug"`
	Plan            string                    `json:"plan"`
	SubscriptionPlan string                   `json:"subscription_plan"` // alias of plan
	LogoURL         *string                   `json:"logo_url,omitempty"`
	Industry        *string                   `json:"industry,omitempty"`
	CompanySize     *string                   `json:"company_size,omitempty"`
	Country         *string                   `json:"country,omitempty"`
	Status          string                    `json:"status"`
	General         entity.GeneralSettings    `json:"general"`
	Automation      entity.AutomationSettings `json:"automation"`
	UpdatedAt       time.Time                 `json:"updated_at"`
}

// UpdateSettingsRequest is a partial update.
type UpdateSettingsRequest struct {
	Name        *string                    `json:"name" validate:"omitempty,min=1,max=200"`
	LogoURL     *string                    `json:"logo_url"`
	Industry    *string                    `json:"industry" validate:"omitempty,max=120"`
	CompanySize *string                    `json:"company_size" validate:"omitempty,max=40"`
	Country     *string                    `json:"country" validate:"omitempty,max=80"`
	Status      *string                    `json:"status" validate:"omitempty,oneof=active suspended trial inactive"`
	General     *entity.GeneralSettings    `json:"general"`
	Automation  *entity.AutomationSettings `json:"automation"`
}
