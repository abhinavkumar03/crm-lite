package notify

import "go.uber.org/zap"

// WhatsAppConfig carries the primitive settings needed to select a WhatsApp
// provider, keeping the notify package free of a dependency on the config
// package (dependency direction points inward).
type WhatsAppConfig struct {
	Provider string // "simulation" | "meta" | "twilio" | "gupshup" | "interakt" | "360dialog"
	APIURL   string
	Token    string
	PhoneID  string
	// Twilio-specific (Account SID may share Token field usage via From).
	AccountSID string
	AuthToken  string
	FromNumber string
}

// EmailConfig selects an email provider and carries connection details.
type EmailConfig struct {
	Provider string // "simulation" | "smtp" | "ses" | "sendgrid" | "mailgun" | "resend"

	// SMTP
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string
	Encryption   string // none | starttls | tls

	// API providers
	APIKey   string
	APIURL   string
	From     string
	ReplyTo  string
}

// BuildWhatsAppProvider selects a WhatsApp provider from configuration. It falls
// back to the simulation provider whenever live credentials are incomplete.
func BuildWhatsAppProvider(cfg WhatsAppConfig, logger *zap.Logger) Provider {
	switch cfg.Provider {
	case "meta":
		if cfg.Token != "" && cfg.PhoneID != "" {
			return NewMetaCloudProvider(cfg.APIURL, cfg.Token, cfg.PhoneID, logger)
		}
	case "twilio":
		if cfg.AccountSID != "" && cfg.AuthToken != "" && cfg.FromNumber != "" {
			return NewTwilioWhatsAppProvider(cfg.AccountSID, cfg.AuthToken, cfg.FromNumber, logger)
		}
	case "gupshup", "interakt", "360dialog":
		// Stub until dedicated adapters land; prefer explicit simulation naming.
		return NewSimulationProvider(cfg.Provider+"-stub", ChannelWhatsApp, logger)
	}
	return NewSimulationProvider("simulation", ChannelWhatsApp, logger)
}

// BuildEmailProvider selects an email provider from configuration.
func BuildEmailProvider(cfg EmailConfig, logger *zap.Logger) Provider {
	switch cfg.Provider {
	case "smtp":
		if cfg.SMTPHost != "" && cfg.SMTPFrom != "" {
			return NewSMTPProvider(cfg, logger)
		}
	case "resend":
		if cfg.APIKey != "" && (cfg.From != "" || cfg.SMTPFrom != "") {
			return NewResendProvider(cfg, logger)
		}
	case "ses", "sendgrid", "mailgun":
		// Adapters follow the same pattern; fall back until credentials wired.
		if cfg.APIKey != "" {
			return NewResendCompatibleStub(cfg.Provider, cfg, logger)
		}
		return NewSimulationProvider(cfg.Provider+"-stub", ChannelEmail, logger)
	}
	return NewSimulationProvider("simulation", ChannelEmail, logger)
}
