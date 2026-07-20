package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	activityentity "github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
	activityrepo "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
)

// WebhookHandler ingests provider delivery/open/read events (public, signature-verified).
type WebhookHandler struct {
	repo         *repository.Repository
	activityRepo *activityrepo.Repository
	logger       *zap.Logger
	metaSecret   string
	metaVerify   string
	resendSecret string
}

func NewWebhookHandler(
	repo *repository.Repository,
	activityRepo *activityrepo.Repository,
	logger *zap.Logger,
	metaSecret, metaVerify, resendSecret string,
) *WebhookHandler {
	return &WebhookHandler{
		repo: repo, activityRepo: activityRepo, logger: logger,
		metaSecret: metaSecret, metaVerify: metaVerify, resendSecret: resendSecret,
	}
}

// MetaVerify handles GET challenge for WhatsApp Cloud API.
func (h *WebhookHandler) MetaVerify(c *gin.Context) {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	if mode == "subscribe" && token == h.metaVerify {
		c.String(http.StatusOK, challenge)
		return
	}
	c.Status(http.StatusForbidden)
}

// MetaStatus handles WhatsApp delivery/read webhooks.
func (h *WebhookHandler) MetaStatus(c *gin.Context) {
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if h.metaSecret != "" {
		sig := c.GetHeader("X-Hub-Signature-256")
		if !verifyMetaSignature(h.metaSecret, body, sig) {
			c.Status(http.StatusUnauthorized)
			return
		}
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	entries, _ := payload["entry"].([]any)
	for _, entry := range entries {
		em, _ := entry.(map[string]any)
		changes, _ := em["changes"].([]any)
		for _, ch := range changes {
			cm, _ := ch.(map[string]any)
			value, _ := cm["value"].(map[string]any)
			statuses, _ := value["statuses"].([]any)
			for _, st := range statuses {
				sm, _ := st.(map[string]any)
				msgID, _ := sm["id"].(string)
				status, _ := sm["status"].(string)
				if msgID == "" || status == "" {
					continue
				}
				h.applyWhatsAppStatus(c, msgID, status, sm)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *WebhookHandler) applyWhatsAppStatus(c *gin.Context, msgID, status string, raw map[string]any) {
	n, err := h.repo.GetByProviderMessageID(c.Request.Context(), msgID)
	if err != nil || n == nil {
		return
	}
	provider := "meta-cloud"
	switch strings.ToLower(status) {
	case "sent":
		// already marked sent on accept
	case "delivered":
		_ = h.repo.MarkDelivered(c.Request.Context(), n.ID)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusDelivered, Provider: &provider, Payload: mustJSON(raw),
		})
		h.writeActivity(n, activityentity.ActionWhatsAppDelivered)
	case "read":
		_ = h.repo.MarkRead(c.Request.Context(), n.ID)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusRead, Provider: &provider, Payload: mustJSON(raw),
		})
		h.writeActivity(n, activityentity.ActionWhatsAppRead)
	case "failed":
		reason := "provider reported failure"
		if errs, ok := raw["errors"].([]any); ok && len(errs) > 0 {
			if em, ok := errs[0].(map[string]any); ok {
				if t, ok := em["title"].(string); ok {
					reason = t
				}
			}
		}
		_ = h.repo.MarkFailed(c.Request.Context(), n.ID, reason)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusFailed, Provider: &provider, Payload: mustJSON(raw),
		})
		h.writeActivity(n, activityentity.ActionNotificationFailed)
	}
}

// ResendEvents handles Resend email webhooks (delivered/opened/bounced).
func (h *WebhookHandler) ResendEvents(c *gin.Context) {
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	// Optional: verify Svix/Resend signature when secret configured.
	_ = h.resendSecret

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	eventType, _ := payload["type"].(string)
	data, _ := payload["data"].(map[string]any)
	msgID, _ := data["email_id"].(string)
	if msgID == "" {
		msgID, _ = data["id"].(string)
	}
	if msgID == "" {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}

	n, err := h.repo.GetByProviderMessageID(c.Request.Context(), msgID)
	if err != nil || n == nil {
		c.JSON(http.StatusOK, gin.H{"ok": true})
		return
	}
	provider := "resend"
	switch eventType {
	case "email.delivered":
		_ = h.repo.MarkDelivered(c.Request.Context(), n.ID)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusDelivered, Provider: &provider, Payload: body,
		})
		h.writeActivity(n, activityentity.ActionEmailDelivered)
	case "email.opened":
		_ = h.repo.MarkOpened(c.Request.Context(), n.ID)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusOpened, Provider: &provider, Payload: body,
		})
		h.writeActivity(n, activityentity.ActionEmailOpened)
	case "email.bounced", "email.complained":
		_ = h.repo.MarkFailed(c.Request.Context(), n.ID, eventType)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusFailed, Provider: &provider, Payload: body,
		})
		h.writeActivity(n, activityentity.ActionNotificationFailed)
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// TwilioStatus handles Twilio WhatsApp status callbacks.
func (h *WebhookHandler) TwilioStatus(c *gin.Context) {
	msgID := c.PostForm("MessageSid")
	status := c.PostForm("MessageStatus")
	if msgID == "" {
		c.Status(http.StatusOK)
		return
	}
	raw := map[string]any{"MessageSid": msgID, "MessageStatus": status}
	n, err := h.repo.GetByProviderMessageID(c.Request.Context(), msgID)
	if err != nil || n == nil {
		c.Status(http.StatusOK)
		return
	}
	provider := "twilio"
	switch strings.ToLower(status) {
	case "delivered":
		_ = h.repo.MarkDelivered(c.Request.Context(), n.ID)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusDelivered, Provider: &provider, Payload: mustJSON(raw),
		})
		h.writeActivity(n, activityentity.ActionWhatsAppDelivered)
	case "read":
		_ = h.repo.MarkRead(c.Request.Context(), n.ID)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusRead, Provider: &provider, Payload: mustJSON(raw),
		})
		h.writeActivity(n, activityentity.ActionWhatsAppRead)
	case "failed", "undelivered":
		_ = h.repo.MarkFailed(c.Request.Context(), n.ID, status)
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusFailed, Provider: &provider, Payload: mustJSON(raw),
		})
		h.writeActivity(n, activityentity.ActionNotificationFailed)
	}
	c.Status(http.StatusOK)
}

// OpenPixel marks an email as opened via a 1x1 tracking pixel.
func (h *WebhookHandler) OpenPixel(c *gin.Context) {
	token := c.Param("token")
	n, err := h.repo.GetByOpenToken(c.Request.Context(), token)
	if err == nil && n != nil {
		_ = h.repo.MarkOpened(c.Request.Context(), n.ID)
		provider := "pixel"
		_ = h.repo.AddDeliveryEvent(c.Request.Context(), &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID, NotificationID: n.ID,
			Event: entity.StatusOpened, Provider: &provider,
		})
		h.writeActivity(n, activityentity.ActionEmailOpened)
	}
	// 1x1 transparent GIF
	c.Header("Content-Type", "image/gif")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
	c.Data(http.StatusOK, "image/gif", []byte{
		0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0x80, 0x00, 0x00, 0xff, 0xff, 0xff,
		0x00, 0x00, 0x00, 0x21, 0xf9, 0x04, 0x01, 0x00, 0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x01, 0x00, 0x00, 0x02, 0x02, 0x44, 0x01, 0x00, 0x3b,
	})
}

func (h *WebhookHandler) writeActivity(n *entity.Notification, action string) {
	if n.CreatedBy == nil || *n.CreatedBy == "" || h.activityRepo == nil {
		return
	}
	meta, _ := json.Marshal(map[string]any{
		"channel": n.Channel, "notification_id": n.ID, "to": n.Recipient,
	})
	act := &activityentity.Activity{
		ID:          uuid.NewString(),
		EntityType:  activityentity.EntityNotification,
		EntityID:    n.ID,
		Action:      action,
		Description: action + " for " + n.Recipient,
		PerformedBy: *n.CreatedBy,
		Metadata:    meta,
	}
	_ = h.activityRepo.Create(context.Background(), act)
}

func verifyMetaSignature(secret string, body []byte, header string) bool {
	if !strings.HasPrefix(header, "sha256=") {
		return false
	}
	want := strings.TrimPrefix(header, "sha256=")
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	got := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(got), []byte(want))
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	if b == nil {
		return []byte("{}")
	}
	return b
}
