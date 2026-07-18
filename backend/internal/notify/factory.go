package notify

import "go.uber.org/zap"

// WhatsAppConfig carries the primitive settings needed to select a WhatsApp
// provider, keeping the notify package free of a dependency on the config
// package (dependency direction points inward).
type WhatsAppConfig struct {
	Provider string // "simulation" | "meta"
	APIURL   string
	Token    string
	PhoneID  string
}

// BuildWhatsAppProvider selects a WhatsApp provider from configuration. It falls
// back to the simulation provider whenever "meta" is not explicitly requested or
// its credentials are incomplete, so the app is always functional out of the box
// and never attempts real delivery without full configuration.
func BuildWhatsAppProvider(cfg WhatsAppConfig, logger *zap.Logger) Provider {
	if cfg.Provider == "meta" && cfg.Token != "" && cfg.PhoneID != "" {
		return NewMetaCloudProvider(cfg.APIURL, cfg.Token, cfg.PhoneID, logger)
	}
	return NewSimulationProvider("simulation", ChannelWhatsApp, logger)
}
