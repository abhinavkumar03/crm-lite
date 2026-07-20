package entity

import "time"

// Notification lifecycle states.
const (
	StatusDraft      = "draft"
	StatusScheduled  = "scheduled"
	StatusQueued     = "queued"
	StatusProcessing = "processing"
	StatusSent       = "sent"
	StatusDelivered  = "delivered"
	StatusOpened     = "opened"
	StatusRead       = "read"
	StatusFailed     = "failed"
	StatusRetrying   = "retrying"
	StatusCancelled  = "cancelled"
)

const (
	ChannelEmail    = "email"
	ChannelWhatsApp = "whatsapp"
)

// Notification is a durable outbound-message record. Content is rendered at
// creation / enqueue time; the worker transitions status and stamps providers.
type Notification struct {
	ID               string
	OrganizationID   string
	Channel          string
	Recipient        string
	CC               []string
	BCC              []string
	Subject          *string
	Body             string
	BodyHTML         *string
	Template         *string
	TemplateID       *string
	Data             []byte // JSONB
	VariablesUsed    []byte // JSONB
	Status           string
	Provider           *string
	ProviderID         *string
	ProviderMessageID  *string
	FromAddress        *string
	ReplyTo            *string
	Error              *string
	LastError          *string
	ProviderResponse   []byte // JSONB
	EntityType         *string
	EntityID           *string
	ModuleID           *string
	AttachmentIDs      []string
	RetryCount         int
	MaxRetries         int
	NextRetryAt        *time.Time
	CreatedBy          *string
	ScheduledAt        *time.Time
	CancelledAt        *time.Time
	QueuedAt           *time.Time
	ProcessingAt       *time.Time
	SentAt             *time.Time
	DeliveredAt        *time.Time
	OpenedAt           *time.Time
	ReadAt             *time.Time
	OpenTrackingToken  *string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
