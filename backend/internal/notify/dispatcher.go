package notify

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// Dispatcher routes messages to the provider registered for their channel.
// It is safe to configure at boot and read concurrently thereafter; providers
// should be registered during application wiring before Dispatch is called.
//
// For per-org providers (M1+), use ResolveAndSend with an OrgProviderResolver
// instead of the process-wide map.
type Dispatcher struct {
	providers map[Channel]Provider
	logger    *zap.Logger
	resolver  OrgProviderResolver
}

// OrgProviderResolver builds a Provider for a given org + channel at send time.
type OrgProviderResolver interface {
	Resolve(ctx context.Context, orgID string, channel Channel) (Provider, error)
}

func NewDispatcher(logger *zap.Logger) *Dispatcher {
	return &Dispatcher{
		providers: make(map[Channel]Provider),
		logger:    logger,
	}
}

// SetResolver enables per-org provider lookup. When set, DispatchWithOrg uses it
// and falls back to registered defaults when Resolve returns nil/error.
func (d *Dispatcher) SetResolver(r OrgProviderResolver) {
	d.resolver = r
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
func (d *Dispatcher) Dispatch(ctx context.Context, msg Message) (SendResult, error) {
	provider, ok := d.providers[msg.Channel]
	if !ok {
		return SendResult{}, fmt.Errorf("notify: no provider registered for channel %q", msg.Channel)
	}
	return provider.Send(ctx, msg)
}

// DispatchWithOrg prefers the org-scoped resolver, then falls back to defaults.
func (d *Dispatcher) DispatchWithOrg(ctx context.Context, orgID string, msg Message) (SendResult, Provider, error) {
	if d.resolver != nil && orgID != "" {
		if p, err := d.resolver.Resolve(ctx, orgID, msg.Channel); err == nil && p != nil {
			result, sendErr := p.Send(ctx, msg)
			return result, p, sendErr
		}
	}
	provider, ok := d.providers[msg.Channel]
	if !ok {
		return SendResult{}, nil, fmt.Errorf("notify: no provider registered for channel %q", msg.Channel)
	}
	result, err := provider.Send(ctx, msg)
	return result, provider, err
}

// ProviderName returns the name of the provider registered for a channel (for
// audit/logging), or "unknown" if none is registered.
func (d *Dispatcher) ProviderName(channel Channel) string {
	if p, ok := d.providers[channel]; ok {
		return p.Name()
	}
	return "unknown"
}
