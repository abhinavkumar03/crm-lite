package notify

import (
	"context"

	"go.uber.org/zap"
)

// SimulationProvider is the default provider for every channel. It performs no
// real network delivery; it structurally logs the message so the full
// notification pipeline (enqueue -> worker -> dispatch -> provider) can be
// exercised end-to-end without external credentials. Real providers added in
// later phases satisfy the same Provider interface and can be swapped in via
// Dispatcher.Register with no changes to callers.
type SimulationProvider struct {
	name    string
	channel Channel
	logger  *zap.Logger
}

func NewSimulationProvider(name string, channel Channel, logger *zap.Logger) *SimulationProvider {
	return &SimulationProvider{
		name:    name,
		channel: channel,
		logger:  logger,
	}
}

func (p *SimulationProvider) Name() string { return p.name }

func (p *SimulationProvider) Channel() Channel { return p.channel }

func (p *SimulationProvider) Send(_ context.Context, msg Message) error {
	if p.logger != nil {
		p.logger.Info("notify: simulated delivery",
			zap.String("provider", p.name),
			zap.String("channel", string(p.channel)),
			zap.String("to", msg.To),
			zap.String("subject", msg.Subject),
			zap.String("template", msg.Template),
		)
	}
	return nil
}
