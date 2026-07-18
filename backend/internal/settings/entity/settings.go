package entity

import "time"

// Organization is the persisted tenant row. Name and the free-form Settings
// JSONB are user-editable from the Settings Center; Slug and Plan are managed
// elsewhere (signup / billing) and are read-only here.
type Organization struct {
	ID        string
	Name      string
	Slug      string
	Plan      string
	Settings  []byte // raw JSONB
	UpdatedAt time.Time
}

// GeneralSettings are org-wide display/locale preferences. They are stored under
// the "general" key of organizations.settings.
type GeneralSettings struct {
	Timezone   string `json:"timezone"`
	DateFormat string `json:"date_format"`
	TimeFormat string `json:"time_format"` // 12h | 24h
	Currency   string `json:"currency"`
	Locale     string `json:"locale"`
	WeekStart  string `json:"week_start"` // sunday | monday
}

// AutomationSettings are org-wide automation preferences (the actual provider
// credentials live in environment config; these are behavioural toggles).
// Stored under the "automation" key of organizations.settings.
type AutomationSettings struct {
	NotificationsEnabled bool   `json:"notifications_enabled"`
	DefaultChannel       string `json:"default_channel"` // whatsapp | email
	DailyDigest          bool   `json:"daily_digest"`
}

// DefaultGeneral returns the baseline general settings used when an organization
// has never saved any.
func DefaultGeneral() GeneralSettings {
	return GeneralSettings{
		Timezone:   "UTC",
		DateFormat: "YYYY-MM-DD",
		TimeFormat: "24h",
		Currency:   "USD",
		Locale:     "en-US",
		WeekStart:  "monday",
	}
}

// DefaultAutomation returns the baseline automation settings.
func DefaultAutomation() AutomationSettings {
	return AutomationSettings{
		NotificationsEnabled: true,
		DefaultChannel:       "whatsapp",
		DailyDigest:          false,
	}
}
