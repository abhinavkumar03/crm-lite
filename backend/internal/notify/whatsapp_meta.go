package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// MetaCloudProvider delivers WhatsApp messages through the Meta (WhatsApp Cloud)
// Graph API. It satisfies the same Provider interface as the simulation
// provider, so it can be swapped in purely via configuration with no changes to
// callers or the worker.
type MetaCloudProvider struct {
	apiURL  string
	token   string
	phoneID string
	client  *http.Client
	logger  *zap.Logger
}

// NewMetaCloudProvider builds a Meta Cloud API provider. apiURL should be the
// Graph API base (e.g. https://graph.facebook.com/v20.0).
func NewMetaCloudProvider(apiURL, token, phoneID string, logger *zap.Logger) *MetaCloudProvider {
	return &MetaCloudProvider{
		apiURL:  strings.TrimRight(apiURL, "/"),
		token:   token,
		phoneID: phoneID,
		client:  &http.Client{Timeout: 15 * time.Second},
		logger:  logger,
	}
}

func (p *MetaCloudProvider) Name() string { return "meta-cloud" }

func (p *MetaCloudProvider) Channel() Channel { return ChannelWhatsApp }

func (p *MetaCloudProvider) Send(ctx context.Context, msg Message) error {
	endpoint := fmt.Sprintf("%s/%s/messages", p.apiURL, p.phoneID)

	// Text message body. Template messages would use a different payload shape;
	// this covers the common session-message case.
	payload := map[string]any{
		"messaging_product": "whatsapp",
		"to":                msg.To,
		"type":              "text",
		"text": map[string]any{
			"preview_url": false,
			"body":        msg.Body,
		},
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify(meta): marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("notify(meta): build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("notify(meta): request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("notify(meta): unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if p.logger != nil {
		p.logger.Info("notify: whatsapp delivered via meta cloud",
			zap.String("to", msg.To),
			zap.String("template", msg.Template),
		)
	}
	return nil
}
