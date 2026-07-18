package processor

import (
	"context"
	"encoding/json"
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
// provider pipeline, transitions the notification's status, and writes an audit
// activity. It implements jobs.NotificationProcessor and runs in the worker.
type Processor struct {
	repo         *repository.Repository
	dispatcher   *notify.Dispatcher
	activityRepo *activityrepo.Repository
	logger       *zap.Logger
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

// Process is idempotent: an already-sent notification is a no-op, and a missing
// one is skipped (returning nil so asynq does not retry forever).
func (p *Processor) Process(ctx context.Context, orgID, id string) error {
	n, err := p.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return err
	}
	if n == nil {
		p.logger.Warn("notification: not found; skipping", zap.String("id", id))
		return nil
	}
	if n.Status == entity.StatusSent {
		return nil
	}

	data := map[string]any{}
	if len(n.Data) > 0 {
		_ = json.Unmarshal(n.Data, &data)
	}

	msg := notify.Message{
		Channel:  notify.Channel(n.Channel),
		To:       n.Recipient,
		Subject:  deref(n.Subject),
		Body:     n.Body,
		Template: deref(n.Template),
		Data:     data,
	}
	provider := p.dispatcher.ProviderName(msg.Channel)

	if err := p.dispatcher.Dispatch(ctx, msg); err != nil {
		if mErr := p.repo.MarkFailed(ctx, n.ID, err.Error()); mErr != nil {
			p.logger.Error("notification: mark failed", zap.Error(mErr))
		}
		p.writeActivity(n, provider, activityentity.ActionNotificationFailed, err.Error())
		// Returning the error lets asynq retry the delivery.
		return err
	}

	if err := p.repo.MarkSent(ctx, n.ID, provider); err != nil {
		return err
	}
	p.writeActivity(n, provider, sentAction(n.Channel), "")
	return nil
}

func sentAction(channel string) string {
	if channel == string(notify.ChannelWhatsApp) {
		return activityentity.ActionWhatsAppSent
	}
	return activityentity.ActionEmailSent
}

// writeActivity records an audit entry. It requires a performing user (the
// activities table demands a non-null performed_by), so notifications created by
// a system path without a user are simply not audited here. Failures are logged,
// never propagated — auditing must not fail delivery.
func (p *Processor) writeActivity(n *entity.Notification, provider, action, failReason string) {
	if n.CreatedBy == nil || *n.CreatedBy == "" {
		return
	}

	meta := map[string]any{
		"channel":  n.Channel,
		"provider": provider,
		"to":       n.Recipient,
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
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
