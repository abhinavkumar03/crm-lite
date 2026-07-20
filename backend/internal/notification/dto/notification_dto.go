package dto

import "time"

// ComposeRequest creates a draft, sends immediately, or schedules delivery.
// mode defaults to "send" for backward compatibility with SendNotificationRequest clients.
type ComposeRequest struct {
	Mode           string         `json:"mode" validate:"omitempty,oneof=draft send schedule"`
	Channel        string         `json:"channel" validate:"required,oneof=email whatsapp"`
	To             string         `json:"to" validate:"required,max=255"`
	CC             []string       `json:"cc"`
	BCC            []string       `json:"bcc"`
	Subject        string         `json:"subject" validate:"omitempty,max=255"`
	Body           string         `json:"body" validate:"required"`
	BodyHTML       string         `json:"body_html"`
	Template       string         `json:"template" validate:"omitempty,max=120"`
	TemplateID     string         `json:"template_id" validate:"omitempty,uuid"`
	Data           map[string]any `json:"data"`
	EntityType     string         `json:"entity_type" validate:"omitempty,max=40"`
	EntityID       string         `json:"entity_id" validate:"omitempty,uuid"`
	ModuleID       string         `json:"module_id" validate:"omitempty,uuid"`
	AttachmentIDs  []string       `json:"attachment_ids"`
	ScheduledAt    *time.Time     `json:"scheduled_at"`
	MaxRetries     *int           `json:"max_retries" validate:"omitempty,min=0,max=10"`
}

// SendNotificationRequest is kept for older clients; maps to ComposeRequest mode=send.
type SendNotificationRequest = ComposeRequest

// NotificationResponse is the API representation of a notification.
type NotificationResponse struct {
	ID               string         `json:"id"`
	Channel          string         `json:"channel"`
	Recipient        string         `json:"recipient"`
	CC               []string       `json:"cc"`
	BCC              []string       `json:"bcc"`
	Subject          *string        `json:"subject"`
	Body             string         `json:"body"`
	BodyHTML         *string        `json:"body_html,omitempty"`
	Template         *string        `json:"template"`
	TemplateID       *string        `json:"template_id,omitempty"`
	Data             map[string]any `json:"data"`
	VariablesUsed    map[string]any `json:"variables_used,omitempty"`
	Status           string         `json:"status"`
	Provider         *string        `json:"provider"`
	Error            *string        `json:"error"`
	LastError        *string        `json:"last_error,omitempty"`
	ProviderResponse map[string]any `json:"provider_response,omitempty"`
	EntityType       *string        `json:"entity_type"`
	EntityID         *string        `json:"entity_id"`
	ModuleID         *string        `json:"module_id,omitempty"`
	AttachmentIDs    []string       `json:"attachment_ids,omitempty"`
	RetryCount       int            `json:"retry_count"`
	MaxRetries       int            `json:"max_retries"`
	CreatedBy        *string        `json:"created_by"`
	ScheduledAt      *time.Time     `json:"scheduled_at,omitempty"`
	CancelledAt      *time.Time     `json:"cancelled_at,omitempty"`
	QueuedAt         *time.Time     `json:"queued_at,omitempty"`
	ProcessingAt     *time.Time     `json:"processing_at,omitempty"`
	SentAt           *time.Time     `json:"sent_at"`
	DeliveredAt      *time.Time     `json:"delivered_at,omitempty"`
	OpenedAt         *time.Time     `json:"opened_at,omitempty"`
	ReadAt           *time.Time     `json:"read_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	Events           []DeliveryEventResponse `json:"events,omitempty"`
}

// DeliveryEventResponse is one lifecycle audit entry.
type DeliveryEventResponse struct {
	ID        string         `json:"id"`
	Event     string         `json:"event"`
	Provider  *string        `json:"provider,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

// ListQuery is the parsed, sanitized set of list parameters.
type ListQuery struct {
	Page       int
	PageSize   int
	Status     string
	Channel    string
	Q          string
	ModuleID   string
	EntityID   string
	EntityType string
	DateFrom   string
	DateTo     string
	TemplateID string
}

// ListResult is a paginated collection of notifications.
type ListResult struct {
	Notifications []NotificationResponse `json:"notifications"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	Total         int                    `json:"total"`
	TotalPages    int                    `json:"total_pages"`
}

// MetricsResponse aggregates delivery stats for the notification center / dashboard.
type MetricsResponse struct {
	EmailsSentToday     int64   `json:"emails_sent_today"`
	WhatsAppSentToday   int64   `json:"whatsapp_sent_today"`
	FailedCount         int64   `json:"failed_count"`
	ScheduledCount      int64   `json:"scheduled_count"`
	DraftCount          int64   `json:"draft_count"`
	TotalSent           int64   `json:"total_sent"`
	TotalDelivered      int64   `json:"total_delivered"`
	DeliveryRate        float64 `json:"delivery_rate"`
	OpenRate            float64 `json:"open_rate"`
	ReadRate            float64 `json:"read_rate"`
}

// CreateTemplateRequest creates an org-scoped reusable template.
type CreateTemplateRequest struct {
	Channel   string   `json:"channel" validate:"required,oneof=email whatsapp"`
	Name      string   `json:"name" validate:"required,min=1,max=120"`
	Category  string   `json:"category" validate:"omitempty,oneof=sales follow_up welcome proposal invoice reminder quotation marketing support custom"`
	Subject   string   `json:"subject" validate:"omitempty,max=255"`
	Body      string   `json:"body" validate:"required"`
	BodyHTML  string   `json:"body_html"`
	Variables []string `json:"variables"`
	IsActive  *bool    `json:"is_active"`
	Status    string   `json:"status" validate:"omitempty,oneof=draft published"`
}

// UpdateTemplateRequest partially updates a template.
type UpdateTemplateRequest struct {
	Name      *string  `json:"name" validate:"omitempty,min=1,max=120"`
	Category  *string  `json:"category" validate:"omitempty,oneof=sales follow_up welcome proposal invoice reminder quotation marketing support custom"`
	Subject   *string  `json:"subject" validate:"omitempty,max=255"`
	Body      *string  `json:"body"`
	BodyHTML  *string  `json:"body_html"`
	Variables []string `json:"variables"`
	IsActive  *bool    `json:"is_active"`
	Status    *string  `json:"status" validate:"omitempty,oneof=draft published"`
}

// TemplateResponse is the API representation of a template.
type TemplateResponse struct {
	ID        string    `json:"id"`
	Channel   string    `json:"channel"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Subject   *string   `json:"subject"`
	Body      string    `json:"body"`
	BodyHTML  *string   `json:"body_html,omitempty"`
	Variables []string  `json:"variables"`
	IsActive  bool      `json:"is_active"`
	Status    string    `json:"status"`
	Version   int       `json:"version"`
	CreatedBy *string   `json:"created_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TemplateListResult is a paginated template collection.
type TemplateListResult struct {
	Templates  []TemplateResponse `json:"templates"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	Total      int                `json:"total"`
	TotalPages int                `json:"total_pages"`
}

// PreviewTemplateRequest renders a template with optional merge data.
type PreviewTemplateRequest struct {
	Data       map[string]any `json:"data"`
	ModuleID   string         `json:"module_id" validate:"omitempty,uuid"`
	EntityID   string         `json:"entity_id" validate:"omitempty,uuid"`
	EntityType string         `json:"entity_type" validate:"omitempty,max=40"`
}

// PreviewTemplateResponse is the rendered preview output.
type PreviewTemplateResponse struct {
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	BodyHTML string `json:"body_html,omitempty"`
}
