package entity

import "time"

// Provider is an org-scoped communication transport configuration.
type Provider struct {
	ID                string
	OrganizationID    string
	Channel           string
	ProviderType      string
	Name              string
	Config            []byte // JSONB non-secret
	SecretsEncrypted  []byte
	IsDefault         bool
	IsActive          bool
	LastHealthAt      *time.Time
	LastError         *string
	CreatedBy         *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	SecretsConfigured bool // computed for API responses
}

// SenderIdentity is a from/reply-to identity for outbound messages.
type SenderIdentity struct {
	ID             string
	OrganizationID string
	ProviderID     *string
	Channel        string
	DisplayName    *string
	FromAddress    string
	ReplyTo        *string
	IsDefault      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
