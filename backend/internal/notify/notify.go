// Package notify defines a channel-agnostic notification pipeline. Concrete
// providers (email/SMTP, WhatsApp via Meta Cloud API, Twilio, Gupshup,
// 360Dialog, ...) implement the Provider interface and are registered against a
// Dispatcher. Callers construct a Message and hand it to the Dispatcher, which
// routes it to the provider registered for that channel (Strategy pattern).
//
// This decouples business logic and the job worker from any specific vendor and
// lets email and WhatsApp share a single delivery pipeline.
package notify

import "context"

// Channel identifies the delivery medium for a notification.
type Channel string

const (
	ChannelEmail    Channel = "email"
	ChannelWhatsApp Channel = "whatsapp"
)

// AttachmentRef points at an already-uploaded file (Cloudinary / attachments table).
type AttachmentRef struct {
	URL      string
	FileName string
	MimeType string
}

// Message is a provider-agnostic notification payload. Template + Data support
// server-rendered templates, while Subject/Body allow pre-rendered content.
type Message struct {
	Channel        Channel
	To             string
	CC             []string
	BCC            []string
	From           string
	ReplyTo        string
	Subject        string
	Body           string
	HTMLBody       string
	Template       string
	// WhatsAppTemplateName / Language / Components support Meta-approved templates.
	WhatsAppTemplateName string
	WhatsAppLanguage     string
	WhatsAppComponents   []map[string]any
	Data                 map[string]any
	Attachments          []AttachmentRef
	IdempotencyKey       string
}

// SendResult is returned after a provider accepts (or rejects) a message.
// ProviderMessageID is required for webhook correlation. Accepted means the
// vendor acknowledged the request; delivery/open/read arrive via webhooks.
type SendResult struct {
	ProviderMessageID string
	RawResponse       map[string]any
	// Simulated is true only for the local SimulationProvider.
	Simulated bool
	// AutoDelivered lets simulation (and only simulation) jump to delivered.
	AutoDelivered bool
}

// Provider is implemented by every concrete notification vendor.
type Provider interface {
	// Name is a human-readable identifier for logging/observability
	// (e.g. "simulation", "meta-cloud", "twilio", "smtp", "resend").
	Name() string
	// Channel reports which medium this provider delivers on.
	Channel() Channel
	// Send delivers the message and returns a result describing acceptance.
	Send(ctx context.Context, msg Message) (SendResult, error)
}
