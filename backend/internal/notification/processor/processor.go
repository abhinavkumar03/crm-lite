package processor

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	activityentity "github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
	activityrepo "github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
)

// Processor delivers a persisted notification: it dispatches through the shared
// provider pipeline, transitions the notification's status, and writes audit
// activities (notification entity + RECORD timeline when linked).
type Processor struct {
	repo         *repository.Repository
	dispatcher   *notify.Dispatcher
	activityRepo *activityrepo.Repository
	logger       *zap.Logger
	publicBase   string
}

func New(
	repo *repository.Repository,
	dispatcher *notify.Dispatcher,
	activityRepo *activityrepo.Repository,
	logger *zap.Logger,
) *Processor {
	return &Processor{
		repo:         repo,
		dispatcher:   dispatcher,
		activityRepo: activityRepo,
		logger:       logger,
	}
}

func (p *Processor) SetPublicBaseURL(base string) {
	p.publicBase = strings.TrimRight(base, "/")
}

// Process is idempotent: an already-sent notification is a no-op, and a missing
// one is skipped (returning nil so asynq does not retry forever).
func (p *Processor) Process(ctx context.Context, orgID, id string) error {
	existing, err := p.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return err
	}
	if existing == nil {
		p.logger.Warn("notification: not found; skipping", zap.String("id", id))
		return nil
	}
	switch existing.Status {
	case entity.StatusSent, entity.StatusDelivered, entity.StatusOpened, entity.StatusRead, entity.StatusCancelled, entity.StatusDraft:
		return nil
	case entity.StatusScheduled:
		if existing.ScheduledAt != nil && !existing.ScheduledAt.After(time.Now().UTC()) {
			if err := p.repo.PromoteScheduled(ctx, id); err != nil {
				return err
			}
		} else {
			return nil
		}
	}

	n, err := p.repo.ClaimForProcessing(ctx, orgID, id)
	if err != nil {
		return err
	}
	if n == nil {
		return nil
	}

	_ = p.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
		OrganizationID: orgID,
		NotificationID: n.ID,
		Event:          entity.StatusProcessing,
	})

	data := map[string]any{}
	if len(n.Data) > 0 {
		_ = json.Unmarshal(n.Data, &data)
	}

	attachments := p.loadAttachments(ctx, n.AttachmentIDs)

	htmlBody := deref(n.BodyHTML)
	if n.Channel == entity.ChannelEmail && n.OpenTrackingToken != nil && *n.OpenTrackingToken != "" {
		base := p.publicBase
		if base == "" {
			base = "http://localhost:8080"
		}
		pixel := `<img src="` + base + `/t/o/` + *n.OpenTrackingToken + `" width="1" height="1" alt="" style="display:none" />`
		if htmlBody == "" {
			htmlBody = "<pre>" + n.Body + "</pre>" + pixel
		} else {
			htmlBody = htmlBody + pixel
		}
	}

	msg := notify.Message{
		Channel:        notify.Channel(n.Channel),
		To:             n.Recipient,
		CC:             n.CC,
		BCC:            n.BCC,
		From:           deref(n.FromAddress),
		ReplyTo:        deref(n.ReplyTo),
		Subject:        deref(n.Subject),
		Body:           n.Body,
		HTMLBody:       htmlBody,
		Template:       deref(n.Template),
		Data:           data,
		Attachments:    attachments,
		IdempotencyKey: n.ID,
	}
	if name, _ := data["whatsapp_template_name"].(string); name != "" {
		msg.WhatsAppTemplateName = name
	}
	if lang, _ := data["whatsapp_language"].(string); lang != "" {
		msg.WhatsAppLanguage = lang
	}

	result, provider, err := p.dispatcher.DispatchWithOrg(ctx, orgID, msg)
	providerName := "unknown"
	if provider != nil {
		providerName = provider.Name()
	}

	if err != nil {
		if n.RetryCount+1 < n.MaxRetries {
			_ = p.repo.MarkRetrying(ctx, n.ID, err.Error())
			_ = p.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
				OrganizationID: orgID,
				NotificationID: n.ID,
				Event:          entity.StatusRetrying,
				Provider:       &providerName,
				Payload:        mustJSON(map[string]any{"error": err.Error(), "retry_count": n.RetryCount + 1}),
			})
			// Do not write failure activity on intermediate retries.
			return err // asynq retries
		}
		if mErr := p.repo.MarkFailed(ctx, n.ID, err.Error()); mErr != nil {
			p.logger.Error("notification: mark failed", zap.Error(mErr))
		}
		_ = p.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
			OrganizationID: orgID,
			NotificationID: n.ID,
			Event:          entity.StatusFailed,
			Provider:       &providerName,
			Payload:        mustJSON(map[string]any{"error": err.Error()}),
		})
		p.writeActivities(n, providerName, activityentity.ActionNotificationFailed, err.Error())
		return err
	}

	response := result.RawResponse
	if response == nil {
		response = map[string]any{}
	}
	response["simulated"] = result.Simulated
	if result.ProviderMessageID != "" {
		response["provider_message_id"] = result.ProviderMessageID
	}

	if err := p.repo.MarkSent(ctx, n.ID, providerName, result.ProviderMessageID, response); err != nil {
		return err
	}
	_ = p.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
		OrganizationID: orgID,
		NotificationID: n.ID,
		Event:          entity.StatusSent,
		Provider:       &providerName,
		Payload:        mustJSON(map[string]any{"provider_message_id": result.ProviderMessageID}),
	})

	// Only simulation may auto-advance to delivered (local demos without webhooks).
	if result.AutoDelivered || result.Simulated {
		_ = p.repo.MarkDelivered(ctx, n.ID)
		_ = p.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
			OrganizationID: orgID,
			NotificationID: n.ID,
			Event:          entity.StatusDelivered,
			Provider:       &providerName,
			Payload:        mustJSON(map[string]any{"auto": true, "simulated": result.Simulated}),
		})
	}

	p.writeActivities(n, providerName, sentAction(n.Channel), "")
	return nil
}

func (p *Processor) loadAttachments(ctx context.Context, ids []string) []notify.AttachmentRef {
	if len(ids) == 0 {
		return nil
	}
	rows, err := p.repo.GetAttachmentsByIDs(ctx, ids)
	if err != nil {
		p.logger.Warn("notification: load attachments", zap.Error(err))
		return nil
	}
	out := make([]notify.AttachmentRef, 0, len(rows))
	for _, a := range rows {
		mime := "application/octet-stream"
		if strings.Contains(strings.ToLower(a.FileName), ".pdf") {
			mime = "application/pdf"
		} else if strings.HasSuffix(strings.ToLower(a.FileName), ".png") {
			mime = "image/png"
		} else if strings.HasSuffix(strings.ToLower(a.FileName), ".jpg") || strings.HasSuffix(strings.ToLower(a.FileName), ".jpeg") {
			mime = "image/jpeg"
		}
		out = append(out, notify.AttachmentRef{
			URL:      a.FileURL,
			FileName: a.FileName,
			MimeType: mime,
		})
	}
	return out
}

// ProcessDueScheduled promotes due scheduled notifications and delivers them.
func (p *Processor) ProcessDueScheduled(ctx context.Context) error {
	due, err := p.repo.ListDueScheduled(ctx, 100)
	if err != nil {
		return err
	}
	for i := range due {
		n := &due[i]
		if err := p.repo.PromoteScheduled(ctx, n.ID); err != nil {
			p.logger.Error("notification: promote scheduled", zap.Error(err), zap.String("id", n.ID))
			continue
		}
		_ = p.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID,
			NotificationID: n.ID,
			Event:          entity.StatusQueued,
			Payload:        mustJSON(map[string]any{"from": "scheduled_tick"}),
		})
		if err := p.Process(ctx, n.OrganizationID, n.ID); err != nil {
			p.logger.Warn("notification: due send failed", zap.Error(err), zap.String("id", n.ID))
		}
	}

	// Also process retrying rows whose next_retry_at is due.
	retries, err := p.repo.ListDueRetries(ctx, 50)
	if err != nil {
		p.logger.Warn("notification: list due retries", zap.Error(err))
	} else {
		for i := range retries {
			n := &retries[i]
			if err := p.repo.MarkQueued(ctx, n.ID); err != nil {
				continue
			}
			if err := p.Process(ctx, n.OrganizationID, n.ID); err != nil {
				p.logger.Warn("notification: retry send failed", zap.Error(err), zap.String("id", n.ID))
			}
		}
	}
	return nil
}

func sentAction(channel string) string {
	if channel == string(notify.ChannelWhatsApp) {
		return activityentity.ActionWhatsAppSent
	}
	return activityentity.ActionEmailSent
}

func (p *Processor) writeActivities(n *entity.Notification, provider, action, failReason string) {
	if n.CreatedBy == nil || *n.CreatedBy == "" {
		return
	}

	meta := map[string]any{
		"channel":         n.Channel,
		"provider":        provider,
		"to":              n.Recipient,
		"notification_id": n.ID,
	}
	if failReason != "" {
		meta["error"] = failReason
	}
	metaBytes, _ := json.Marshal(meta)
	description := "Notification " + action + " to " + n.Recipient

	act := &activityentity.Activity{
		ID:          uuid.NewString(),
		EntityType:  activityentity.EntityNotification,
		EntityID:    n.ID,
		Action:      action,
		Description: description,
		PerformedBy: *n.CreatedBy,
		Metadata:    metaBytes,
		CreatedAt:   time.Now().UTC(),
	}
	if err := p.activityRepo.Create(context.Background(), act); err != nil {
		p.logger.Error("notification: write activity", zap.Error(err))
	}

	if n.ModuleID != nil && *n.ModuleID != "" && n.EntityID != nil && *n.EntityID != "" {
		if err := p.repo.CreateRecordActivity(
			context.Background(),
			n.OrganizationID, *n.ModuleID, *n.EntityID, *n.CreatedBy,
			action, description, metaBytes,
		); err != nil {
			p.logger.Error("notification: write record activity", zap.Error(err))
		}
	}
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	if b == nil {
		return []byte("{}")
	}
	return b
}
