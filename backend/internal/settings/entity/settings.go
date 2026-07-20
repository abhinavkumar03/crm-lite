package entity

import "time"

// Organization is the persisted tenant row.
type Organization struct {
	ID          string
	Name        string
	Slug        string
	Plan        string
	LogoURL     *string
	Description *string
	Industry    *string
	CompanySize *string
	Country     *string
	Status      string
	CreatedBy   *string
	Settings    []byte // raw JSONB
	UpdatedAt   time.Time
}

// GeneralSettings are org-wide display/locale preferences under settings.general.
type GeneralSettings struct {
	Timezone   string `json:"timezone"`
	DateFormat string `json:"date_format"`
	TimeFormat string `json:"time_format"` // 12h | 24h
	Currency   string `json:"currency"`
	Locale     string `json:"locale"`
	WeekStart  string `json:"week_start"` // sunday | monday
}

// AutomationSettings are org-wide automation preferences under settings.automation.
type AutomationSettings struct {
	NotificationsEnabled bool   `json:"notifications_enabled"`
	DefaultChannel       string `json:"default_channel"` // whatsapp | email
	DailyDigest          bool   `json:"daily_digest"`
}

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

func DefaultAutomation() AutomationSettings {
	return AutomationSettings{
		NotificationsEnabled: true,
		DefaultChannel:       "whatsapp",
		DailyDigest:          false,
	}
}
