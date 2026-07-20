package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// TwilioWhatsAppProvider sends WhatsApp messages via Twilio Messaging API.
type TwilioWhatsAppProvider struct {
	accountSID string
	authToken  string
	from       string
	client     *http.Client
	logger     *zap.Logger
}

func NewTwilioWhatsAppProvider(accountSID, authToken, from string, logger *zap.Logger) *TwilioWhatsAppProvider {
	from = strings.TrimSpace(from)
	if from != "" && !strings.HasPrefix(from, "whatsapp:") {
		from = "whatsapp:" + from
	}
	return &TwilioWhatsAppProvider{
		accountSID: accountSID,
		authToken:  authToken,
		from:       from,
		client:     &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}
}

func (p *TwilioWhatsAppProvider) Name() string    { return "twilio" }
func (p *TwilioWhatsAppProvider) Channel() Channel { return ChannelWhatsApp }

func (p *TwilioWhatsAppProvider) Send(ctx context.Context, msg Message) (SendResult, error) {
	to := strings.TrimSpace(msg.To)
	if to != "" && !strings.HasPrefix(to, "whatsapp:") {
		to = "whatsapp:" + to
	}

	form := url.Values{}
	form.Set("To", to)
	form.Set("From", p.from)
	if msg.WhatsAppTemplateName != "" {
		// Content SID / template body — Twilio Content API uses ContentSid.
		form.Set("ContentSid", msg.WhatsAppTemplateName)
		if len(msg.Data) > 0 {
			if raw, err := json.Marshal(msg.Data); err == nil {
				form.Set("ContentVariables", string(raw))
			}
		}
	} else {
		form.Set("Body", msg.Body)
	}

	endpoint := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", p.accountSID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return SendResult{}, err
	}
	req.SetBasicAuth(p.accountSID, p.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.client.Do(req)
	if err != nil {
		return SendResult{}, fmt.Errorf("notify(twilio): request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
	parsed := map[string]any{}
	_ = json.Unmarshal(body, &parsed)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return SendResult{RawResponse: parsed}, fmt.Errorf("notify(twilio): status %d: %s", resp.StatusCode, string(body))
	}

	msgID, _ := parsed["sid"].(string)
	if p.logger != nil {
		p.logger.Info("notify: whatsapp accepted via twilio",
			zap.String("to", msg.To),
			zap.String("provider_message_id", msgID),
		)
	}
	return SendResult{
		ProviderMessageID: msgID,
		RawResponse:       parsed,
	}, nil
}
