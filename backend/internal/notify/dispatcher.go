package notify

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Dispatcher routes messages to the provider registered for their channel.
// It is safe to configure at boot and read concurrently thereafter; providers
// should be registered during application wiring before Dispatch is called.
type Dispatcher struct {
	providers map[Channel]Provider
	logger    *zap.Logger
}

func NewDispatcher(logger *zap.Logger) *Dispatcher {
	return &Dispatcher{
		providers: make(map[Channel]Provider),
		logger:    logger,
	}
}

// Register wires a provider for its channel. A later registration for the same
// channel overrides the earlier one, which makes swapping vendors a one-line
// change at composition time.
func (d *Dispatcher) Register(p Provider) {
	d.providers[p.Channel()] = p

	if d.logger != nil {
		d.logger.Info("notify: provider registered",
			zap.String("provider", p.Name()),
			zap.String("channel", string(p.Channel())),
		)
	}
}

// Dispatch delivers msg via the provider registered for its channel.
func (d *Dispatcher) Dispatch(ctx context.Context, msg Message) error {
	provider, ok := d.providers[msg.Channel]
	if !ok {
		return fmt.Errorf("notify: no provider registered for channel %q", msg.Channel)
	}

	return provider.Send(ctx, msg)
}

// ProviderName returns the name of the provider registered for a channel (for
// audit/logging), or "unknown" if none is registered.
func (d *Dispatcher) ProviderName(channel Channel) string {
	if p, ok := d.providers[channel]; ok {
		return p.Name()
	}
	return "unknown"
}
