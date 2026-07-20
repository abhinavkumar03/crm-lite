package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
)

const notificationColumns = `
	id, organization_id, channel, recipient, COALESCE(cc, '{}'), COALESCE(bcc, '{}'),
	subject, body, body_html, template, template_id, data, COALESCE(variables_used, '{}'),
	status, provider, provider_id, provider_message_id, from_address, reply_to,
	error, last_error, COALESCE(provider_response, '{}'),
	entity_type, entity_id, module_id, COALESCE(attachment_ids, '{}'),
	retry_count, max_retries, next_retry_at, created_by,
	scheduled_at, cancelled_at, queued_at, processing_at, sent_at,
	delivered_at, opened_at, read_at, open_tracking_token, created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanNotification(row pgx.Row, n *entity.Notification) error {
	return row.Scan(
		&n.ID, &n.OrganizationID, &n.Channel, &n.Recipient, &n.CC, &n.BCC,
		&n.Subject, &n.Body, &n.BodyHTML, &n.Template, &n.TemplateID, &n.Data, &n.VariablesUsed,
		&n.Status, &n.Provider, &n.ProviderID, &n.ProviderMessageID, &n.FromAddress, &n.ReplyTo,
		&n.Error, &n.LastError, &n.ProviderResponse,
		&n.EntityType, &n.EntityID, &n.ModuleID, &n.AttachmentIDs,
		&n.RetryCount, &n.MaxRetries, &n.NextRetryAt, &n.CreatedBy,
		&n.ScheduledAt, &n.CancelledAt, &n.QueuedAt, &n.ProcessingAt, &n.SentAt,
		&n.DeliveredAt, &n.OpenedAt, &n.ReadAt, &n.OpenTrackingToken, &n.CreatedAt, &n.UpdatedAt,
	)
}

func (r *Repository) Create(ctx context.Context, n *entity.Notification) error {
	if n.CC == nil {
		n.CC = []string{}
	}
	if n.BCC == nil {
		n.BCC = []string{}
	}
	if n.AttachmentIDs == nil {
		n.AttachmentIDs = []string{}
	}
	if n.Data == nil {
		n.Data = []byte("{}")
	}
	if n.VariablesUsed == nil {
		n.VariablesUsed = []byte("{}")
	}
	if n.ProviderResponse == nil {
		n.ProviderResponse = []byte("{}")
	}
	if n.MaxRetries <= 0 {
		n.MaxRetries = 3
	}

	return r.db.QueryRow(ctx, `
		INSERT INTO notifications (
			organization_id, channel, recipient, cc, bcc, subject, body, body_html,
			template, template_id, data, variables_used, status, entity_type, entity_id,
			module_id, attachment_ids, retry_count, max_retries, created_by,
			scheduled_at, queued_at, from_address, reply_to, open_tracking_token
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25)
		RETURNING id, status, created_at, updated_at
	`,
		n.OrganizationID, n.Channel, n.Recipient, n.CC, n.BCC, n.Subject, n.Body, n.BodyHTML,
		n.Template, n.TemplateID, n.Data, n.VariablesUsed, n.Status, n.EntityType, n.EntityID,
		n.ModuleID, n.AttachmentIDs, n.RetryCount, n.MaxRetries, n.CreatedBy,
		n.ScheduledAt, n.QueuedAt, n.FromAddress, n.ReplyTo, n.OpenTrackingToken,
	).Scan(&n.ID, &n.Status, &n.CreatedAt, &n.UpdatedAt)
}

func (r *Repository) UpdateDraft(ctx context.Context, n *entity.Notification) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications SET
			channel = $3, recipient = $4, cc = $5, bcc = $6, subject = $7, body = $8,
			body_html = $9, template = $10, template_id = $11, data = $12, variables_used = $13,
			entity_type = $14, entity_id = $15, module_id = $16, attachment_ids = $17,
			scheduled_at = $18, status = $19, queued_at = $20, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND status = 'draft'
	`,
		n.ID, n.OrganizationID, n.Channel, n.Recipient, n.CC, n.BCC, n.Subject, n.Body,
		n.BodyHTML, n.Template, n.TemplateID, n.Data, n.VariablesUsed,
		n.EntityType, n.EntityID, n.ModuleID, n.AttachmentIDs,
		n.ScheduledAt, n.Status, n.QueuedAt,
	)
	return err
}

func (r *Repository) GetByID(ctx context.Context, orgID, id string) (*entity.Notification, error) {
	var n entity.Notification
	err := scanNotification(r.db.QueryRow(ctx, `
		SELECT `+notificationColumns+`
		FROM notifications
		WHERE id = $1 AND organization_id = $2
	`, id, orgID), &n)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) List(ctx context.Context, orgID string, q dto.ListQuery) ([]entity.Notification, int, error) {
	conds := []string{"organization_id = $1"}
	args := []any{orgID}
	next := 2

	if q.Status != "" {
		conds = append(conds, fmt.Sprintf("status = $%d", next))
		args = append(args, q.Status)
		next++
	}
	if q.Channel != "" {
		conds = append(conds, fmt.Sprintf("channel = $%d", next))
		args = append(args, q.Channel)
		next++
	}
	if q.ModuleID != "" {
		conds = append(conds, fmt.Sprintf("module_id = $%d", next))
		args = append(args, q.ModuleID)
		next++
	}
	if q.EntityID != "" {
		conds = append(conds, fmt.Sprintf("entity_id = $%d", next))
		args = append(args, q.EntityID)
		next++
	}
	if q.EntityType != "" {
		conds = append(conds, fmt.Sprintf("entity_type = $%d", next))
		args = append(args, q.EntityType)
		next++
	}
	if q.TemplateID != "" {
		conds = append(conds, fmt.Sprintf("template_id = $%d", next))
		args = append(args, q.TemplateID)
		next++
	}
	if q.Q != "" {
		conds = append(conds, fmt.Sprintf(
			"(recipient ILIKE $%d OR COALESCE(subject,'') ILIKE $%d OR COALESCE(template,'') ILIKE $%d)",
			next, next, next,
		))
		args = append(args, "%"+q.Q+"%")
		next++
	}
	if q.DateFrom != "" {
		conds = append(conds, fmt.Sprintf("created_at >= $%d::timestamptz", next))
		args = append(args, q.DateFrom)
		next++
	}
	if q.DateTo != "" {
		conds = append(conds, fmt.Sprintf("created_at <= $%d::timestamptz", next))
		args = append(args, q.DateTo)
		next++
	}

	where := "WHERE " + strings.Join(conds, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notifications "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitPH := fmt.Sprintf("$%d", next)
	offsetPH := fmt.Sprintf("$%d", next+1)
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.Query(ctx,
		"SELECT "+notificationColumns+" FROM notifications "+where+
			" ORDER BY created_at DESC LIMIT "+limitPH+" OFFSET "+offsetPH, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]entity.Notification, 0)
	for rows.Next() {
		var n entity.Notification
		if err := scanNotification(rows, &n); err != nil {
			return nil, 0, err
		}
		items = append(items, n)
	}
	return items, total, rows.Err()
}

// ClaimForProcessing atomically moves queued/retrying → processing.
func (r *Repository) ClaimForProcessing(ctx context.Context, orgID, id string) (*entity.Notification, error) {
	var n entity.Notification
	err := scanNotification(r.db.QueryRow(ctx, `
		UPDATE notifications
		SET status = 'processing', processing_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND status IN ('queued', 'retrying')
		RETURNING `+notificationColumns+`
	`, id, orgID), &n)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) MarkSent(ctx context.Context, id, provider, providerMessageID string, response map[string]any) error {
	payload, _ := json.Marshal(response)
	if payload == nil {
		payload = []byte("{}")
	}
	var msgID any
	if providerMessageID != "" {
		msgID = providerMessageID
	}
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'sent', provider = $2, provider_message_id = COALESCE($3, provider_message_id),
		    error = NULL, last_error = NULL,
		    provider_response = $4, sent_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id, provider, msgID, payload)
	return err
}

func (r *Repository) MarkDelivered(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'delivered', delivered_at = COALESCE(delivered_at, NOW()), updated_at = NOW()
		WHERE id = $1 AND status IN ('sent', 'processing', 'queued', 'delivered')
	`, id)
	return err
}

func (r *Repository) MarkOpened(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'opened', opened_at = COALESCE(opened_at, NOW()),
		    delivered_at = COALESCE(delivered_at, NOW()), updated_at = NOW()
		WHERE id = $1 AND status IN ('sent', 'delivered', 'opened')
	`, id)
	return err
}

func (r *Repository) MarkRead(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'read', read_at = COALESCE(read_at, NOW()),
		    delivered_at = COALESCE(delivered_at, NOW()), updated_at = NOW()
		WHERE id = $1 AND status IN ('sent', 'delivered', 'opened', 'read')
	`, id)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, id, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'failed', error = $2, last_error = $2, updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	return err
}

func (r *Repository) MarkRetrying(ctx context.Context, id, reason string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'retrying', retry_count = retry_count + 1,
		    error = $2, last_error = $2,
		    next_retry_at = NOW() + (LEAST(POWER(2, retry_count + 1), 60) || ' minutes')::interval,
		    updated_at = NOW()
		WHERE id = $1
	`, id, reason)
	return err
}

func (r *Repository) GetByProviderMessageID(ctx context.Context, providerMessageID string) (*entity.Notification, error) {
	var n entity.Notification
	err := scanNotification(r.db.QueryRow(ctx, `
		SELECT `+notificationColumns+`
		FROM notifications
		WHERE provider_message_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, providerMessageID), &n)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) GetByOpenToken(ctx context.Context, token string) (*entity.Notification, error) {
	var n entity.Notification
	err := scanNotification(r.db.QueryRow(ctx, `
		SELECT `+notificationColumns+`
		FROM notifications
		WHERE open_tracking_token = $1
		LIMIT 1
	`, token), &n)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) ListDueRetries(ctx context.Context, limit int) ([]entity.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT `+notificationColumns+`
		FROM notifications
		WHERE status = 'retrying' AND next_retry_at IS NOT NULL AND next_retry_at <= NOW()
		ORDER BY next_retry_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]entity.Notification, 0)
	for rows.Next() {
		var n entity.Notification
		if err := scanNotification(rows, &n); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return items, rows.Err()
}

// AttachmentMeta is a minimal attachment projection for outbound sends.
type AttachmentMeta struct {
	ID       string
	FileName string
	FileURL  string
}

func (r *Repository) GetAttachmentsByIDs(ctx context.Context, ids []string) ([]AttachmentMeta, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := r.db.Query(ctx, `
		SELECT id, file_name, file_url FROM attachments WHERE id = ANY($1::uuid[])
	`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]AttachmentMeta, 0)
	for rows.Next() {
		var a AttachmentMeta
		if err := rows.Scan(&a.ID, &a.FileName, &a.FileURL); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repository) LinkAttachments(ctx context.Context, notificationID string, attachmentIDs []string) error {
	for _, id := range attachmentIDs {
		if id == "" {
			continue
		}
		_, err := r.db.Exec(ctx, `
			INSERT INTO notification_attachments (notification_id, attachment_id)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, notificationID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) MarkQueued(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'queued', queued_at = NOW(), updated_at = NOW()
		WHERE id = $1
	`, id)
	return err
}

func (r *Repository) CancelScheduled(ctx context.Context, orgID, id string) (*entity.Notification, error) {
	var n entity.Notification
	err := scanNotification(r.db.QueryRow(ctx, `
		UPDATE notifications
		SET status = 'cancelled', cancelled_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND organization_id = $2 AND status = 'scheduled'
		RETURNING `+notificationColumns+`
	`, id, orgID), &n)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *Repository) ListDueScheduled(ctx context.Context, limit int) ([]entity.Notification, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT `+notificationColumns+`
		FROM notifications
		WHERE status = 'scheduled' AND scheduled_at IS NOT NULL AND scheduled_at <= NOW()
		ORDER BY scheduled_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Notification, 0)
	for rows.Next() {
		var n entity.Notification
		if err := scanNotification(rows, &n); err != nil {
			return nil, err
		}
		items = append(items, n)
	}
	return items, rows.Err()
}

func (r *Repository) PromoteScheduled(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notifications
		SET status = 'queued', queued_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND status = 'scheduled'
	`, id)
	return err
}

func (r *Repository) AddDeliveryEvent(ctx context.Context, e *entity.DeliveryEvent) error {
	if e.ID == "" {
		e.ID = uuid.NewString()
	}
	if e.Payload == nil {
		e.Payload = []byte("{}")
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO notification_delivery_events (
			id, organization_id, notification_id, event, provider, payload
		) VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING created_at
	`, e.ID, e.OrganizationID, e.NotificationID, e.Event, e.Provider, e.Payload).Scan(&e.CreatedAt)
}

func (r *Repository) ListDeliveryEvents(ctx context.Context, orgID, notificationID string) ([]entity.DeliveryEvent, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, organization_id, notification_id, event, provider, payload, created_at
		FROM notification_delivery_events
		WHERE organization_id = $1 AND notification_id = $2
		ORDER BY created_at ASC
	`, orgID, notificationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.DeliveryEvent, 0)
	for rows.Next() {
		var e entity.DeliveryEvent
		if err := rows.Scan(&e.ID, &e.OrganizationID, &e.NotificationID, &e.Event, &e.Provider, &e.Payload, &e.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, e)
	}
	return items, rows.Err()
}

func (r *Repository) CreateRecordActivity(
	ctx context.Context,
	orgID, moduleID, recordID, userID, action, description string,
	metadata json.RawMessage,
) error {
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO activities (
			id, entity_type, entity_id, action, description, performed_by,
			metadata, organization_id, module_id, created_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$8,NOW())
	`, uuid.NewString(), recordID, action, description, userID, metadata, orgID, moduleID)
	return err
}

func (r *Repository) Metrics(ctx context.Context, orgID string) (*dto.MetricsResponse, error) {
	m := &dto.MetricsResponse{}
	var emailOpened, whatsappRead, emailSent, whatsappSent int64
	err := r.db.QueryRow(ctx, `
		SELECT
			COUNT(*) FILTER (
				WHERE channel = 'email'
				  AND status IN ('sent','delivered','opened','read')
				  AND COALESCE(sent_at, created_at)::date = CURRENT_DATE
			),
			COUNT(*) FILTER (
				WHERE channel = 'whatsapp'
				  AND status IN ('sent','delivered','opened','read')
				  AND COALESCE(sent_at, created_at)::date = CURRENT_DATE
			),
			COUNT(*) FILTER (WHERE status IN ('failed','retrying')),
			COUNT(*) FILTER (WHERE status = 'scheduled'),
			COUNT(*) FILTER (WHERE status = 'draft'),
			COUNT(*) FILTER (WHERE status IN ('sent','delivered','opened','read')),
			COUNT(*) FILTER (WHERE status IN ('delivered','opened','read')),
			COUNT(*) FILTER (WHERE status IN ('opened','read') AND channel = 'email'),
			COUNT(*) FILTER (WHERE status = 'read' AND channel = 'whatsapp'),
			COUNT(*) FILTER (WHERE channel = 'email' AND status IN ('sent','delivered','opened','read')),
			COUNT(*) FILTER (WHERE channel = 'whatsapp' AND status IN ('sent','delivered','opened','read'))
		FROM notifications
		WHERE organization_id = $1
	`, orgID).Scan(
		&m.EmailsSentToday, &m.WhatsAppSentToday, &m.FailedCount, &m.ScheduledCount, &m.DraftCount,
		&m.TotalSent, &m.TotalDelivered,
		&emailOpened, &whatsappRead, &emailSent, &whatsappSent,
	)
	if err != nil {
		return nil, err
	}

	if m.TotalSent > 0 {
		m.DeliveryRate = float64(m.TotalDelivered) / float64(m.TotalSent) * 100
	}
	if emailSent > 0 {
		m.OpenRate = float64(emailOpened) / float64(emailSent) * 100
	}
	if whatsappSent > 0 {
		m.ReadRate = float64(whatsappRead) / float64(whatsappSent) * 100
	}
	return m, nil
}

// --- Templates ---

const templateColumns = `
	id, organization_id, channel, name, category, subject, body, body_html,
	variables, is_active, COALESCE(status, 'published'), COALESCE(version, 1),
	created_by, created_at, updated_at
`

func scanTemplate(row pgx.Row, t *entity.Template) error {
	return row.Scan(
		&t.ID, &t.OrganizationID, &t.Channel, &t.Name, &t.Category, &t.Subject,
		&t.Body, &t.BodyHTML, &t.Variables, &t.IsActive, &t.Status, &t.Version,
		&t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
}

func (r *Repository) CreateTemplate(ctx context.Context, t *entity.Template) error {
	if t.Variables == nil {
		t.Variables = []byte("[]")
	}
	if t.Status == "" {
		t.Status = entity.TemplateStatusPublished
	}
	if t.Version <= 0 {
		t.Version = 1
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO notification_templates (
			organization_id, channel, name, category, subject, body, body_html,
			variables, is_active, status, version, created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING id, created_at, updated_at
	`,
		t.OrganizationID, t.Channel, t.Name, t.Category, t.Subject, t.Body, t.BodyHTML,
		t.Variables, t.IsActive, t.Status, t.Version, t.CreatedBy,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *Repository) GetTemplate(ctx context.Context, orgID, id string) (*entity.Template, error) {
	var t entity.Template
	err := scanTemplate(r.db.QueryRow(ctx, `
		SELECT `+templateColumns+` FROM notification_templates
		WHERE id = $1 AND organization_id = $2
	`, id, orgID), &t)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) UpdateTemplate(ctx context.Context, t *entity.Template) error {
	_, err := r.db.Exec(ctx, `
		UPDATE notification_templates SET
			name = $3, category = $4, subject = $5, body = $6, body_html = $7,
			variables = $8, is_active = $9, status = $10, version = $11, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`, t.ID, t.OrganizationID, t.Name, t.Category, t.Subject, t.Body, t.BodyHTML, t.Variables, t.IsActive, t.Status, t.Version)
	return err
}

func (r *Repository) PublishTemplate(ctx context.Context, orgID, id string) (*entity.Template, error) {
	var t entity.Template
	err := scanTemplate(r.db.QueryRow(ctx, `
		UPDATE notification_templates
		SET status = 'published', version = version + 1, is_active = TRUE, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
		RETURNING `+templateColumns+`
	`, id, orgID), &t)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) DeleteTemplate(ctx context.Context, orgID, id string) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM notification_templates WHERE id = $1 AND organization_id = $2
	`, id, orgID)
	return err
}

func (r *Repository) ListTemplates(ctx context.Context, orgID, channel, category string, page, pageSize int) ([]entity.Template, int, error) {
	conds := []string{"organization_id = $1"}
	args := []any{orgID}
	next := 2
	if channel != "" {
		conds = append(conds, fmt.Sprintf("channel = $%d", next))
		args = append(args, channel)
		next++
	}
	if category != "" {
		conds = append(conds, fmt.Sprintf("category = $%d", next))
		args = append(args, category)
		next++
	}
	where := "WHERE " + strings.Join(conds, " AND ")

	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM notification_templates "+where, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	limitPH := fmt.Sprintf("$%d", next)
	offsetPH := fmt.Sprintf("$%d", next+1)
	args = append(args, pageSize, (page-1)*pageSize)

	rows, err := r.db.Query(ctx,
		"SELECT "+templateColumns+" FROM notification_templates "+where+
			" ORDER BY updated_at DESC LIMIT "+limitPH+" OFFSET "+offsetPH, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]entity.Template, 0)
	for rows.Next() {
		var t entity.Template
		if err := scanTemplate(rows, &t); err != nil {
			return nil, 0, err
		}
		items = append(items, t)
	}
	return items, total, rows.Err()
}

// LoadMergeContext builds dotted merge variables from a CRM record + org.
func (r *Repository) LoadMergeContext(ctx context.Context, orgID, moduleID, entityID, userID string) (map[string]any, error) {
	out := map[string]any{
		"today":        time.Now().UTC().Format("2006-01-02"),
		"current_date": time.Now().UTC().Format("2006-01-02"),
	}

	var orgName string
	_ = r.db.QueryRow(ctx, `SELECT name FROM organizations WHERE id = $1`, orgID).Scan(&orgName)
	if orgName != "" {
		out["workspace"] = map[string]any{"name": orgName}
		out["workspace.name"] = orgName
	}

	if userID != "" {
		var fullName, email string
		_ = r.db.QueryRow(ctx, `
			SELECT COALESCE(NULLIF(TRIM(name), ''), email), email
			FROM users WHERE id = $1
		`, userID).Scan(&fullName, &email)
		out["owner"] = map[string]any{"name": fullName, "email": email}
		out["owner.name"] = fullName
		out["owner.email"] = email
	}

	if moduleID == "" || entityID == "" {
		return out, nil
	}

	var apiName string
	var data []byte
	err := r.db.QueryRow(ctx, `
		SELECT m.api_name, r.data
		FROM records r
		JOIN modules m ON m.id = r.module_id
		WHERE r.id = $1 AND r.organization_id = $2 AND r.module_id = $3
	`, entityID, orgID, moduleID).Scan(&apiName, &data)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return out, nil
		}
		return out, err
	}

	recordData := map[string]any{}
	if len(data) > 0 {
		_ = json.Unmarshal(data, &recordData)
	}
	moduleKey := strings.TrimSuffix(apiName, "s")
	if moduleKey == "" {
		moduleKey = apiName
	}
	// Common aliases for merge templates.
	aliases := []string{moduleKey, apiName, "record", "lead", "contact", "task"}
	seen := map[string]bool{}
	for _, key := range aliases {
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		out[key] = recordData
		for field, val := range recordData {
			out[key+"."+field] = val
		}
	}
	return out, nil
}
