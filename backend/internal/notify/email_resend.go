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

// ResendProvider sends email via the Resend HTTP API.
type ResendProvider struct {
	apiKey  string
	apiURL  string
	from    string
	replyTo string
	client  *http.Client
	logger  *zap.Logger
}

func NewResendProvider(cfg EmailConfig, logger *zap.Logger) *ResendProvider {
	apiURL := cfg.APIURL
	if apiURL == "" {
		apiURL = "https://api.resend.com"
	}
	from := firstNonEmpty(cfg.From, cfg.SMTPFrom)
	return &ResendProvider{
		apiKey:  cfg.APIKey,
		apiURL:  strings.TrimRight(apiURL, "/"),
		from:    from,
		replyTo: cfg.ReplyTo,
		client:  &http.Client{Timeout: 30 * time.Second},
		logger:  logger,
	}
}

func (p *ResendProvider) Name() string    { return "resend" }
func (p *ResendProvider) Channel() Channel { return ChannelEmail }

func (p *ResendProvider) Send(ctx context.Context, msg Message) (SendResult, error) {
	from := firstNonEmpty(msg.From, p.from)
	if from == "" {
		return SendResult{}, fmt.Errorf("notify(resend): from address required")
	}

	payload := map[string]any{
		"from":    from,
		"to":      []string{msg.To},
		"subject": msg.Subject,
	}
	if len(msg.CC) > 0 {
		payload["cc"] = msg.CC
	}
	if len(msg.BCC) > 0 {
		payload["bcc"] = msg.BCC
	}
	replyTo := firstNonEmpty(msg.ReplyTo, p.replyTo)
	if replyTo != "" {
		payload["reply_to"] = replyTo
	}
	if msg.HTMLBody != "" {
		payload["html"] = msg.HTMLBody
		if msg.Body != "" {
			payload["text"] = msg.Body
		}
	} else {
		payload["text"] = msg.Body
	}
	if len(msg.Attachments) > 0 {
		atts := make([]map[string]any, 0, len(msg.Attachments))
		for _, a := range msg.Attachments {
			atts = append(atts, map[string]any{
				"filename": a.FileName,
				"path":     a.URL,
			})
		}
		payload["attachments"] = atts
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return SendResult{}, fmt.Errorf("notify(resend): marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.apiURL+"/emails", bytes.NewReader(raw))
	if err != nil {
		return SendResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return SendResult{}, fmt.Errorf("notify(resend): request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	parsed := map[string]any{}
	_ = json.Unmarshal(body, &parsed)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SendResult{RawResponse: parsed}, fmt.Errorf("notify(resend): status %d: %s", resp.StatusCode, string(body))
	}

	msgID, _ := parsed["id"].(string)
	if p.logger != nil {
		p.logger.Info("notify: email accepted via resend",
			zap.String("to", msg.To),
			zap.String("provider_message_id", msgID),
		)
	}
	return SendResult{
		ProviderMessageID: msgID,
		RawResponse:       parsed,
	}, nil
}

// ResendCompatibleStub is a temporary bridge for SES/SendGrid/Mailgun until
// dedicated adapters are complete — uses Resend-shaped config when only API key exists.
type ResendCompatibleStub struct {
	name string
	*ResendProvider
}

func NewResendCompatibleStub(name string, cfg EmailConfig, logger *zap.Logger) Provider {
	p := NewResendProvider(cfg, logger)
	return &ResendCompatibleStub{name: name, ResendProvider: p}
}

func (p *ResendCompatibleStub) Name() string { return p.name }
