package dto

import "time"

// SendNotificationRequest queues a notification for delivery. Subject/Body may
// contain {{placeholder}} tokens resolved from Data at send time.
type SendNotificationRequest struct {
	Channel    string         `json:"channel" validate:"required,oneof=email whatsapp"`
	To         string         `json:"to" validate:"required,max=255"`
	Subject    string         `json:"subject" validate:"omitempty,max=255"`
	Body       string         `json:"body" validate:"required"`
	Template   string         `json:"template" validate:"omitempty,max=120"`
	Data       map[string]any `json:"data"`
	EntityType string         `json:"entity_type" validate:"omitempty,max=40"`
	EntityID   string         `json:"entity_id" validate:"omitempty,uuid"`
}

// NotificationResponse is the API representation of a notification.
type NotificationResponse struct {
	ID         string         `json:"id"`
	Channel    string         `json:"channel"`
	Recipient  string         `json:"recipient"`
	Subject    *string        `json:"subject"`
	Body       string         `json:"body"`
	Template   *string        `json:"template"`
	Data       map[string]any `json:"data"`
	Status     string         `json:"status"`
	Provider   *string        `json:"provider"`
	Error      *string        `json:"error"`
	EntityType *string        `json:"entity_type"`
	EntityID   *string        `json:"entity_id"`
	CreatedBy  *string        `json:"created_by"`
	SentAt     *time.Time     `json:"sent_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// ListQuery is the parsed, sanitized set of list parameters.
type ListQuery struct {
	Page     int
	PageSize int
	Status   string
	Channel  string
}

// ListResult is a paginated collection of notifications.
type ListResult struct {
	Notifications []NotificationResponse `json:"notifications"`
	Page          int                    `json:"page"`
	PageSize      int                    `json:"page_size"`
	Total         int                    `json:"total"`
	TotalPages    int                    `json:"total_pages"`
}
