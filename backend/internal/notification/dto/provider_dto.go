package dto

import "time"

type ProviderUpsertRequest struct {
	Channel      string         `json:"channel" binding:"required"`
	ProviderType string         `json:"provider_type" binding:"required"`
	Name         string         `json:"name" binding:"required"`
	Config       map[string]any `json:"config"`
	Secrets      map[string]any `json:"secrets"`
	IsDefault    *bool          `json:"is_default"`
	IsActive     *bool          `json:"is_active"`
}

type ProviderResponse struct {
	ID                string         `json:"id"`
	Channel           string         `json:"channel"`
	ProviderType      string         `json:"provider_type"`
	Name              string         `json:"name"`
	Config            map[string]any `json:"config"`
	SecretsConfigured bool           `json:"secrets_configured"`
	IsDefault         bool           `json:"is_default"`
	IsActive          bool           `json:"is_active"`
	LastHealthAt      *time.Time     `json:"last_health_at,omitempty"`
	LastError         *string        `json:"last_error,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

type ProviderTestRequest struct {
	To string `json:"to" binding:"required"`
}

type SenderUpsertRequest struct {
	ProviderID  *string `json:"provider_id"`
	Channel     string  `json:"channel" binding:"required"`
	DisplayName *string `json:"display_name"`
	FromAddress string  `json:"from_address" binding:"required"`
	ReplyTo     *string `json:"reply_to"`
	IsDefault   *bool   `json:"is_default"`
}

type SenderResponse struct {
	ID          string    `json:"id"`
	ProviderID  *string   `json:"provider_id,omitempty"`
	Channel     string    `json:"channel"`
	DisplayName *string   `json:"display_name,omitempty"`
	FromAddress string    `json:"from_address"`
	ReplyTo     *string   `json:"reply_to,omitempty"`
	IsDefault   bool      `json:"is_default"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PreferenceResponse struct {
	EmailEnabled    bool `json:"email_enabled"`
	WhatsAppEnabled bool `json:"whatsapp_enabled"`
	Transactional   bool `json:"transactional"`
	Marketing       bool `json:"marketing"`
	DoNotDisturb    bool `json:"do_not_disturb"`
}

type PreferenceUpdateRequest struct {
	EmailEnabled    *bool `json:"email_enabled"`
	WhatsAppEnabled *bool `json:"whatsapp_enabled"`
	Transactional   *bool `json:"transactional"`
	Marketing       *bool `json:"marketing"`
	DoNotDisturb    *bool `json:"do_not_disturb"`
}
