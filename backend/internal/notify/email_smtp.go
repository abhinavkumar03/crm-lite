package notify

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"go.uber.org/zap"
)

// SMTPProvider sends email via SMTP (plain / STARTTLS / TLS).
type SMTPProvider struct {
	host       string
	port       int
	username   string
	password   string
	from       string
	replyTo    string
	encryption string
	logger     *zap.Logger
}

func NewSMTPProvider(cfg EmailConfig, logger *zap.Logger) *SMTPProvider {
	port := cfg.SMTPPort
	if port == 0 {
		port = 587
	}
	from := cfg.SMTPFrom
	if from == "" {
		from = cfg.From
	}
	return &SMTPProvider{
		host:       cfg.SMTPHost,
		port:       port,
		username:   cfg.SMTPUsername,
		password:   cfg.SMTPPassword,
		from:       from,
		replyTo:    firstNonEmpty(cfg.ReplyTo),
		encryption: strings.ToLower(cfg.Encryption),
		logger:     logger,
	}
}

func (p *SMTPProvider) Name() string    { return "smtp" }
func (p *SMTPProvider) Channel() Channel { return ChannelEmail }

func (p *SMTPProvider) Send(ctx context.Context, msg Message) (SendResult, error) {
	from := firstNonEmpty(msg.From, p.from)
	if from == "" {
		return SendResult{}, fmt.Errorf("notify(smtp): from address required")
	}
	replyTo := firstNonEmpty(msg.ReplyTo, p.replyTo)

	var body strings.Builder
	body.WriteString(fmt.Sprintf("From: %s\r\n", from))
	body.WriteString(fmt.Sprintf("To: %s\r\n", msg.To))
	if len(msg.CC) > 0 {
		body.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(msg.CC, ", ")))
	}
	if replyTo != "" {
		body.WriteString(fmt.Sprintf("Reply-To: %s\r\n", replyTo))
	}
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	body.WriteString("MIME-Version: 1.0\r\n")
	if msg.HTMLBody != "" {
		boundary := fmt.Sprintf("crmlite_%d", time.Now().UnixNano())
		body.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary))
		body.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s\r\n", boundary, msg.Body))
		body.WriteString(fmt.Sprintf("--%s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s\r\n", boundary, msg.HTMLBody))
		body.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	} else {
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		body.WriteString(msg.Body)
	}

	addr := fmt.Sprintf("%s:%d", p.host, p.port)
	recipients := append([]string{msg.To}, msg.CC...)
	recipients = append(recipients, msg.BCC...)

	var err error
	switch p.encryption {
	case "tls", "ssl":
		err = p.sendTLS(ctx, addr, from, recipients, []byte(body.String()))
	default:
		err = p.sendStartTLSOrPlain(ctx, addr, from, recipients, []byte(body.String()))
	}
	if err != nil {
		return SendResult{}, err
	}

	msgID := fmt.Sprintf("smtp_%s_%d", msg.IdempotencyKey, time.Now().UnixNano())
	if p.logger != nil {
		p.logger.Info("notify: email accepted via smtp",
			zap.String("to", msg.To),
			zap.String("provider_message_id", msgID),
		)
	}
	return SendResult{
		ProviderMessageID: msgID,
		RawResponse:       map[string]any{"transport": "smtp", "host": p.host},
	}, nil
}

func (p *SMTPProvider) sendStartTLSOrPlain(ctx context.Context, addr, from string, to []string, msg []byte) error {
	d := net.Dialer{Timeout: 15 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("notify(smtp): dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		return fmt.Errorf("notify(smtp): client: %w", err)
	}
	defer client.Close()

	if p.encryption == "" || p.encryption == "starttls" {
		if ok, _ := client.Extension("STARTTLS"); ok {
			cfg := &tls.Config{ServerName: p.host, MinVersion: tls.VersionTLS12}
			if err := client.StartTLS(cfg); err != nil {
				return fmt.Errorf("notify(smtp): starttls: %w", err)
			}
		}
	}
	if p.username != "" {
		auth := smtp.PlainAuth("", p.username, p.password, p.host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("notify(smtp): auth: %w", err)
		}
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("notify(smtp): mail: %w", err)
	}
	for _, rcpt := range to {
		if strings.TrimSpace(rcpt) == "" {
			continue
		}
		if err := client.Rcpt(rcpt); err != nil {
			return fmt.Errorf("notify(smtp): rcpt %s: %w", rcpt, err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("notify(smtp): data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return fmt.Errorf("notify(smtp): write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("notify(smtp): close: %w", err)
	}
	return client.Quit()
}

func (p *SMTPProvider) sendTLS(ctx context.Context, addr, from string, to []string, msg []byte) error {
	d := tls.Dialer{
		Config: &tls.Config{ServerName: p.host, MinVersion: tls.VersionTLS12},
		NetDialer: &net.Dialer{Timeout: 15 * time.Second},
	}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("notify(smtp): tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, p.host)
	if err != nil {
		return fmt.Errorf("notify(smtp): tls client: %w", err)
	}
	defer client.Close()

	if p.username != "" {
		auth := smtp.PlainAuth("", p.username, p.password, p.host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("notify(smtp): auth: %w", err)
		}
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if strings.TrimSpace(rcpt) == "" {
			continue
		}
		if err := client.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		_ = w.Close()
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return client.Quit()
}
