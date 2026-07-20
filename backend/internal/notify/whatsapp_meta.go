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
		client:  &http.Client{Timeout: 30 * time.Second},
		logger:  logger,
	}
}

func (p *MetaCloudProvider) Name() string { return "meta-cloud" }

func (p *MetaCloudProvider) Channel() Channel { return ChannelWhatsApp }

func (p *MetaCloudProvider) Send(ctx context.Context, msg Message) (SendResult, error) {
	endpoint := fmt.Sprintf("%s/%s/messages", p.apiURL, p.phoneID)

	var payload map[string]any
	if msg.WhatsAppTemplateName != "" {
		payload = map[string]any{
			"messaging_product": "whatsapp",
			"to":                msg.To,
			"type":              "template",
			"template": map[string]any{
				"name": msg.WhatsAppTemplateName,
				"language": map[string]any{
					"code": firstNonEmpty(msg.WhatsAppLanguage, "en_US"),
				},
			},
		}
		if len(msg.WhatsAppComponents) > 0 {
			payload["template"].(map[string]any)["components"] = msg.WhatsAppComponents
		}
	} else {
		payload = map[string]any{
			"messaging_product": "whatsapp",
			"to":                msg.To,
			"type":              "text",
			"text": map[string]any{
				"preview_url": false,
				"body":        msg.Body,
			},
		}
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return SendResult{}, fmt.Errorf("notify(meta): marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(raw))
	if err != nil {
		return SendResult{}, fmt.Errorf("notify(meta): build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return SendResult{}, fmt.Errorf("notify(meta): request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	parsed := map[string]any{}
	_ = json.Unmarshal(body, &parsed)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SendResult{RawResponse: parsed}, fmt.Errorf("notify(meta): unexpected status %d: %s", resp.StatusCode, string(body))
	}

	msgID := extractMetaMessageID(parsed)
	if p.logger != nil {
		p.logger.Info("notify: whatsapp accepted via meta cloud",
			zap.String("to", msg.To),
			zap.String("template", msg.Template),
			zap.String("provider_message_id", msgID),
		)
	}
	return SendResult{
		ProviderMessageID: msgID,
		RawResponse:       parsed,
	}, nil
}

func extractMetaMessageID(parsed map[string]any) string {
	messages, ok := parsed["messages"].([]any)
	if !ok || len(messages) == 0 {
		return ""
	}
	first, ok := messages[0].(map[string]any)
	if !ok {
		return ""
	}
	if id, ok := first["id"].(string); ok {
		return id
	}
	return ""
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
