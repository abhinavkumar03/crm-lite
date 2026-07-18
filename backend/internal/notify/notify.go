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

// Message is a provider-agnostic notification payload. Template + Data support
// server-rendered templates, while Subject/Body allow pre-rendered content.
type Message struct {
	Channel  Channel
	To       string
	Subject  string
	Body     string
	Template string
	Data     map[string]any
}

// Provider is implemented by every concrete notification vendor.
type Provider interface {
	// Name is a human-readable identifier for logging/observability
	// (e.g. "simulation", "meta-cloud", "twilio").
	Name() string
	// Channel reports which medium this provider delivers on.
	Channel() Channel
	// Send delivers the message or returns an error describing the failure.
	Send(ctx context.Context, msg Message) error
}
