package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/hibiken/asynq"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notify"
)

const (
	DefaultPageSize = 20
	MaxPageSize     = 100
)

var (
	ErrNotFound       = errors.New("notification not found")
	ErrInvalidState   = errors.New("invalid notification state for this action")
	ErrScheduleRequired = errors.New("scheduled_at is required and must be in the future")
	ErrEmailSubject   = errors.New("subject is required for email")
)

// Repository is the persistence contract for notifications and templates.
type Repository interface {
	Create(ctx context.Context, n *entity.Notification) error
	UpdateDraft(ctx context.Context, n *entity.Notification) error
	GetByID(ctx context.Context, orgID, id string) (*entity.Notification, error)
	List(ctx context.Context, orgID string, q dto.ListQuery) ([]entity.Notification, int, error)
	MarkSent(ctx context.Context, id, provider, providerMessageID string, response map[string]any) error
	MarkFailed(ctx context.Context, id, reason string) error
	MarkRetrying(ctx context.Context, id, reason string) error
	MarkQueued(ctx context.Context, id string) error
	CancelScheduled(ctx context.Context, orgID, id string) (*entity.Notification, error)
	PromoteScheduled(ctx context.Context, id string) error
	ListDueScheduled(ctx context.Context, limit int) ([]entity.Notification, error)
	AddDeliveryEvent(ctx context.Context, e *entity.DeliveryEvent) error
	ListDeliveryEvents(ctx context.Context, orgID, notificationID string) ([]entity.DeliveryEvent, error)
	Metrics(ctx context.Context, orgID string) (*dto.MetricsResponse, error)
	LoadMergeContext(ctx context.Context, orgID, moduleID, entityID, userID string) (map[string]any, error)
	LinkAttachments(ctx context.Context, notificationID string, attachmentIDs []string) error

	CreateTemplate(ctx context.Context, t *entity.Template) error
	GetTemplate(ctx context.Context, orgID, id string) (*entity.Template, error)
	UpdateTemplate(ctx context.Context, t *entity.Template) error
	PublishTemplate(ctx context.Context, orgID, id string) (*entity.Template, error)
	DeleteTemplate(ctx context.Context, orgID, id string) error
	ListTemplates(ctx context.Context, orgID, channel, category string, page, pageSize int) ([]entity.Template, int, error)
}

// Enqueuer publishes jobs onto the async queue (satisfied by *jobs.Producer).
type Enqueuer interface {
	Publish(ctx context.Context, job jobs.Job, opts ...asynq.Option) error
}

type Service struct {
	repo     Repository
	enqueuer Enqueuer
}

func New(repo Repository, enqueuer Enqueuer) *Service {
	return &Service{repo: repo, enqueuer: enqueuer}
}

// Compose creates a draft, queues an immediate send, or schedules for later.
func (s *Service) Compose(ctx context.Context, orgID, userID string, req dto.ComposeRequest) (*dto.NotificationResponse, error) {
	mode := req.Mode
	if mode == "" {
		mode = "send"
	}
	if req.Channel == entity.ChannelEmail && req.Subject == "" && mode != "draft" {
		return nil, ErrEmailSubject
	}
	if mode == "schedule" {
		if req.ScheduledAt == nil || !req.ScheduledAt.After(time.Now().UTC()) {
			return nil, ErrScheduleRequired
		}
	}

	data, variablesUsed, subject, body, bodyHTML, err := s.renderContent(ctx, orgID, userID, req)
	if err != nil {
		return nil, err
	}
	dataBytes, _ := json.Marshal(data)
	varsBytes, _ := json.Marshal(variablesUsed)

	status := entity.StatusQueued
	var queuedAt *time.Time
	now := time.Now().UTC()
	switch mode {
	case "draft":
		status = entity.StatusDraft
	case "schedule":
		status = entity.StatusScheduled
	default:
		queuedAt = &now
	}

	maxRetries := 3
	if req.MaxRetries != nil {
		maxRetries = *req.MaxRetries
	}

	n := &entity.Notification{
		OrganizationID: orgID,
		Channel:        req.Channel,
		Recipient:      req.To,
		CC:             req.CC,
		BCC:            req.BCC,
		Subject:        ptrOrNil(subject),
		Body:           body,
		BodyHTML:       ptrOrNil(bodyHTML),
		Template:       ptrOrNil(req.Template),
		TemplateID:     ptrOrNil(req.TemplateID),
		Data:           dataBytes,
		VariablesUsed:  varsBytes,
		Status:         status,
		EntityType:     ptrOrNil(req.EntityType),
		EntityID:       ptrOrNil(req.EntityID),
		ModuleID:       ptrOrNil(req.ModuleID),
		AttachmentIDs:  req.AttachmentIDs,
		MaxRetries:     maxRetries,
		CreatedBy:      &userID,
		ScheduledAt:    req.ScheduledAt,
		QueuedAt:       queuedAt,
	}
	if req.Channel == entity.ChannelEmail {
		token := newOpenToken()
		n.OpenTrackingToken = &token
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}

	if len(req.AttachmentIDs) > 0 {
		_ = s.repo.LinkAttachments(ctx, n.ID, req.AttachmentIDs)
	}

	_ = s.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
		OrganizationID: orgID,
		NotificationID: n.ID,
		Event:          status,
		Payload:        mustJSON(map[string]any{"mode": mode}),
	})

	switch mode {
	case "send":
		if err := s.enqueueSend(ctx, orgID, userID, n.ID, nil); err != nil {
			_ = s.repo.MarkFailed(ctx, n.ID, "failed to enqueue: "+err.Error())
			return nil, err
		}
	case "schedule":
		if err := s.enqueueSend(ctx, orgID, userID, n.ID, req.ScheduledAt); err != nil {
			_ = s.repo.MarkFailed(ctx, n.ID, "failed to schedule: "+err.Error())
			return nil, err
		}
	}

	resp := toResponse(n)
	return &resp, nil
}

// Send is an alias for Compose with mode=send (backward compatible).
func (s *Service) Send(ctx context.Context, orgID, userID string, req dto.ComposeRequest) (*dto.NotificationResponse, error) {
	req.Mode = "send"
	return s.Compose(ctx, orgID, userID, req)
}

func (s *Service) Retry(ctx context.Context, orgID, userID, id string) (*dto.NotificationResponse, error) {
	n, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, ErrNotFound
	}
	if n.Status != entity.StatusFailed && n.Status != entity.StatusRetrying {
		return nil, ErrInvalidState
	}
	if n.RetryCount >= n.MaxRetries {
		return nil, fmt.Errorf("%w: max retries exceeded", ErrInvalidState)
	}

	if err := s.repo.MarkQueued(ctx, n.ID); err != nil {
		return nil, err
	}
	_ = s.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
		OrganizationID: orgID,
		NotificationID: n.ID,
		Event:          entity.StatusQueued,
		Payload:        mustJSON(map[string]any{"retry": true, "attempt": n.RetryCount + 1}),
	})

	if err := s.enqueueSend(ctx, orgID, userID, n.ID, nil); err != nil {
		_ = s.repo.MarkFailed(ctx, n.ID, "failed to enqueue retry: "+err.Error())
		return nil, err
	}

	updated, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	resp := toResponse(updated)
	return &resp, nil
}

func (s *Service) Cancel(ctx context.Context, orgID, id string) (*dto.NotificationResponse, error) {
	n, err := s.repo.CancelScheduled(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		existing, gerr := s.repo.GetByID(ctx, orgID, id)
		if gerr != nil {
			return nil, gerr
		}
		if existing == nil {
			return nil, ErrNotFound
		}
		return nil, ErrInvalidState
	}
	_ = s.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
		OrganizationID: orgID,
		NotificationID: n.ID,
		Event:          entity.StatusCancelled,
	})
	resp := toResponse(n)
	return &resp, nil
}

func (s *Service) Get(ctx context.Context, orgID, id string) (*dto.NotificationResponse, error) {
	n, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	resp := toResponse(n)
	events, err := s.repo.ListDeliveryEvents(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	resp.Events = make([]dto.DeliveryEventResponse, 0, len(events))
	for i := range events {
		resp.Events = append(resp.Events, toEventResponse(&events[i]))
	}
	return &resp, nil
}

func (s *Service) List(ctx context.Context, orgID string, q dto.ListQuery) (*dto.ListResult, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = DefaultPageSize
	}
	if q.PageSize > MaxPageSize {
		q.PageSize = MaxPageSize
	}

	items, total, err := s.repo.List(ctx, orgID, q)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.NotificationResponse, 0, len(items))
	for i := range items {
		responses = append(responses, toResponse(&items[i]))
	}

	return &dto.ListResult{
		Notifications: responses,
		Page:          q.Page,
		PageSize:      q.PageSize,
		Total:         total,
		TotalPages:    int(math.Max(1, math.Ceil(float64(total)/float64(q.PageSize)))),
	}, nil
}

func (s *Service) Metrics(ctx context.Context, orgID string) (*dto.MetricsResponse, error) {
	return s.repo.Metrics(ctx, orgID)
}

// ProcessDueScheduled promotes due scheduled rows and enqueues them (worker tick).
func (s *Service) ProcessDueScheduled(ctx context.Context) (int, error) {
	due, err := s.repo.ListDueScheduled(ctx, 100)
	if err != nil {
		return 0, err
	}
	count := 0
	for i := range due {
		n := &due[i]
		if err := s.repo.PromoteScheduled(ctx, n.ID); err != nil {
			continue
		}
		userID := ""
		if n.CreatedBy != nil {
			userID = *n.CreatedBy
		}
		if err := s.enqueueSend(ctx, n.OrganizationID, userID, n.ID, nil); err != nil {
			_ = s.repo.MarkFailed(ctx, n.ID, "failed to enqueue scheduled: "+err.Error())
			continue
		}
		_ = s.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
			OrganizationID: n.OrganizationID,
			NotificationID: n.ID,
			Event:          entity.StatusQueued,
			Payload:        mustJSON(map[string]any{"from": "scheduled"}),
		})
		count++
	}
	return count, nil
}

// --- Templates ---

func (s *Service) CreateTemplate(ctx context.Context, orgID, userID string, req dto.CreateTemplateRequest) (*dto.TemplateResponse, error) {
	category := req.Category
	if category == "" {
		category = entity.CategoryCustom
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	vars, _ := json.Marshal(req.Variables)
	if req.Variables == nil {
		vars = []byte("[]")
	}
	t := &entity.Template{
		OrganizationID: orgID,
		Channel:        req.Channel,
		Name:           req.Name,
		Category:       category,
		Subject:        ptrOrNil(req.Subject),
		Body:           req.Body,
		BodyHTML:       ptrOrNil(req.BodyHTML),
		Variables:      vars,
		IsActive:       active,
		Status:         entity.TemplateStatusDraft,
		Version:        1,
		CreatedBy:      &userID,
	}
	if req.Status == entity.TemplateStatusPublished {
		t.Status = entity.TemplateStatusPublished
	}
	if err := s.repo.CreateTemplate(ctx, t); err != nil {
		return nil, err
	}
	resp := toTemplateResponse(t)
	return &resp, nil
}

func (s *Service) GetTemplate(ctx context.Context, orgID, id string) (*dto.TemplateResponse, error) {
	t, err := s.repo.GetTemplate(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	resp := toTemplateResponse(t)
	return &resp, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, orgID, id string, req dto.UpdateTemplateRequest) (*dto.TemplateResponse, error) {
	t, err := s.repo.GetTemplate(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Category != nil {
		t.Category = *req.Category
	}
	if req.Subject != nil {
		t.Subject = req.Subject
	}
	if req.Body != nil {
		t.Body = *req.Body
	}
	if req.BodyHTML != nil {
		t.BodyHTML = req.BodyHTML
	}
	if req.Variables != nil {
		t.Variables, _ = json.Marshal(req.Variables)
	}
	if req.IsActive != nil {
		t.IsActive = *req.IsActive
	}
	if req.Status != nil {
		t.Status = *req.Status
	}
	if err := s.repo.UpdateTemplate(ctx, t); err != nil {
		return nil, err
	}
	resp := toTemplateResponse(t)
	return &resp, nil
}

func (s *Service) PublishTemplate(ctx context.Context, orgID, id string) (*dto.TemplateResponse, error) {
	t, err := s.repo.PublishTemplate(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	resp := toTemplateResponse(t)
	return &resp, nil
}

func (s *Service) UpdateDraft(ctx context.Context, orgID, userID, id string, req dto.ComposeRequest) (*dto.NotificationResponse, error) {
	n, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, ErrNotFound
	}
	if n.Status != entity.StatusDraft {
		return nil, ErrInvalidState
	}
	data, variablesUsed, subject, body, bodyHTML, err := s.renderContent(ctx, orgID, userID, req)
	if err != nil {
		return nil, err
	}
	dataBytes, _ := json.Marshal(data)
	varsBytes, _ := json.Marshal(variablesUsed)
	n.Channel = req.Channel
	n.Recipient = req.To
	n.CC = req.CC
	n.BCC = req.BCC
	n.Subject = ptrOrNil(subject)
	n.Body = body
	n.BodyHTML = ptrOrNil(bodyHTML)
	n.Template = ptrOrNil(req.Template)
	n.TemplateID = ptrOrNil(req.TemplateID)
	n.Data = dataBytes
	n.VariablesUsed = varsBytes
	n.EntityType = ptrOrNil(req.EntityType)
	n.EntityID = ptrOrNil(req.EntityID)
	n.ModuleID = ptrOrNil(req.ModuleID)
	n.AttachmentIDs = req.AttachmentIDs
	n.Status = entity.StatusDraft
	if err := s.repo.UpdateDraft(ctx, n); err != nil {
		return nil, err
	}
	if len(req.AttachmentIDs) > 0 {
		_ = s.repo.LinkAttachments(ctx, n.ID, req.AttachmentIDs)
	}
	resp := toResponse(n)
	return &resp, nil
}

func (s *Service) SendDraft(ctx context.Context, orgID, userID, id string) (*dto.NotificationResponse, error) {
	n, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, ErrNotFound
	}
	if n.Status != entity.StatusDraft {
		return nil, ErrInvalidState
	}
	if err := s.repo.MarkQueued(ctx, n.ID); err != nil {
		return nil, err
	}
	_ = s.repo.AddDeliveryEvent(ctx, &entity.DeliveryEvent{
		OrganizationID: orgID,
		NotificationID: n.ID,
		Event:          entity.StatusQueued,
		Payload:        mustJSON(map[string]any{"from": "draft"}),
	})
	if err := s.enqueueSend(ctx, orgID, userID, n.ID, nil); err != nil {
		_ = s.repo.MarkFailed(ctx, n.ID, "failed to enqueue: "+err.Error())
		return nil, err
	}
	updated, err := s.repo.GetByID(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	resp := toResponse(updated)
	return &resp, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, orgID, id string) error {
	return s.repo.DeleteTemplate(ctx, orgID, id)
}

func (s *Service) ListTemplates(ctx context.Context, orgID, channel, category string, page, pageSize int) (*dto.TemplateListResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	items, total, err := s.repo.ListTemplates(ctx, orgID, channel, category, page, pageSize)
	if err != nil {
		return nil, err
	}
	out := make([]dto.TemplateResponse, 0, len(items))
	for i := range items {
		out = append(out, toTemplateResponse(&items[i]))
	}
	return &dto.TemplateListResult{
		Templates:  out,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: int(math.Max(1, math.Ceil(float64(total)/float64(pageSize)))),
	}, nil
}

func (s *Service) PreviewTemplate(ctx context.Context, orgID, userID, id string, req dto.PreviewTemplateRequest) (*dto.PreviewTemplateResponse, error) {
	t, err := s.repo.GetTemplate(ctx, orgID, id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	data := map[string]any{}
	for k, v := range req.Data {
		data[k] = v
	}
	ctxData, _ := s.repo.LoadMergeContext(ctx, orgID, req.ModuleID, req.EntityID, userID)
	for k, v := range ctxData {
		if _, exists := data[k]; !exists {
			data[k] = v
		}
	}
	subject := ""
	if t.Subject != nil {
		subject = notify.Render(*t.Subject, data)
	}
	bodyHTML := ""
	if t.BodyHTML != nil {
		bodyHTML = notify.Render(*t.BodyHTML, data)
	}
	return &dto.PreviewTemplateResponse{
		Subject:  subject,
		Body:     notify.Render(t.Body, data),
		BodyHTML: bodyHTML,
	}, nil
}

func (s *Service) renderContent(ctx context.Context, orgID, userID string, req dto.ComposeRequest) (data, varsUsed map[string]any, subject, body, bodyHTML string, err error) {
	data = map[string]any{}
	if req.Data != nil {
		for k, v := range req.Data {
			data[k] = v
		}
	}
	ctxData, _ := s.repo.LoadMergeContext(ctx, orgID, req.ModuleID, req.EntityID, userID)
	for k, v := range ctxData {
		if _, exists := data[k]; !exists {
			data[k] = v
		}
	}

	subjectSrc := req.Subject
	bodySrc := req.Body
	bodyHTMLSrc := req.BodyHTML

	if req.TemplateID != "" {
		t, terr := s.repo.GetTemplate(ctx, orgID, req.TemplateID)
		if terr != nil {
			return nil, nil, "", "", "", terr
		}
		if t != nil {
			if subjectSrc == "" && t.Subject != nil {
				subjectSrc = *t.Subject
			}
			if bodySrc == "" {
				bodySrc = t.Body
			}
			if bodyHTMLSrc == "" && t.BodyHTML != nil {
				bodyHTMLSrc = *t.BodyHTML
			}
			if req.Template == "" {
				req.Template = t.Name
			}
		}
	}

	subject = notify.Render(subjectSrc, data)
	body = notify.Render(bodySrc, data)
	bodyHTML = notify.Render(bodyHTMLSrc, data)
	varsUsed = data
	return data, varsUsed, subject, body, bodyHTML, nil
}

func (s *Service) enqueueSend(ctx context.Context, orgID, userID, id string, at *time.Time) error {
	job := jobs.Job{
		Type:   jobs.JobSendNotification,
		UserID: userID,
		Payload: map[string]interface{}{
			"notification_id": id,
			"org_id":          orgID,
		},
	}
	opts := []asynq.Option{}
	if at != nil {
		opts = append(opts, asynq.ProcessAt(*at))
	}
	return s.enqueuer.Publish(ctx, job, opts...)
}

func toResponse(n *entity.Notification) dto.NotificationResponse {
	data := map[string]any{}
	if len(n.Data) > 0 {
		_ = json.Unmarshal(n.Data, &data)
	}
	vars := map[string]any{}
	if len(n.VariablesUsed) > 0 {
		_ = json.Unmarshal(n.VariablesUsed, &vars)
	}
	provResp := map[string]any{}
	if len(n.ProviderResponse) > 0 {
		_ = json.Unmarshal(n.ProviderResponse, &provResp)
	}
	cc := n.CC
	if cc == nil {
		cc = []string{}
	}
	bcc := n.BCC
	if bcc == nil {
		bcc = []string{}
	}
	atts := n.AttachmentIDs
	if atts == nil {
		atts = []string{}
	}
	errMsg := n.Error
	if n.LastError != nil {
		errMsg = n.LastError
	}
	return dto.NotificationResponse{
		ID: n.ID, Channel: n.Channel, Recipient: n.Recipient, CC: cc, BCC: bcc,
		Subject: n.Subject, Body: n.Body, BodyHTML: n.BodyHTML, Template: n.Template,
		TemplateID: n.TemplateID, Data: data, VariablesUsed: vars, Status: n.Status,
		Provider: n.Provider, Error: errMsg, LastError: n.LastError, ProviderResponse: provResp,
		EntityType: n.EntityType, EntityID: n.EntityID, ModuleID: n.ModuleID,
		AttachmentIDs: atts, RetryCount: n.RetryCount, MaxRetries: n.MaxRetries,
		CreatedBy: n.CreatedBy, ScheduledAt: n.ScheduledAt, CancelledAt: n.CancelledAt,
		QueuedAt: n.QueuedAt, ProcessingAt: n.ProcessingAt, SentAt: n.SentAt,
		DeliveredAt: n.DeliveredAt, OpenedAt: n.OpenedAt, ReadAt: n.ReadAt,
		CreatedAt: n.CreatedAt, UpdatedAt: n.UpdatedAt,
	}
}

func toEventResponse(e *entity.DeliveryEvent) dto.DeliveryEventResponse {
	payload := map[string]any{}
	if len(e.Payload) > 0 {
		_ = json.Unmarshal(e.Payload, &payload)
	}
	return dto.DeliveryEventResponse{
		ID: e.ID, Event: e.Event, Provider: e.Provider, Payload: payload, CreatedAt: e.CreatedAt,
	}
}

func toTemplateResponse(t *entity.Template) dto.TemplateResponse {
	vars := []string{}
	if len(t.Variables) > 0 {
		_ = json.Unmarshal(t.Variables, &vars)
	}
	return dto.TemplateResponse{
		ID: t.ID, Channel: t.Channel, Name: t.Name, Category: t.Category,
		Subject: t.Subject, Body: t.Body, BodyHTML: t.BodyHTML, Variables: vars,
		IsActive: t.IsActive, Status: firstNonEmpty(t.Status, entity.TemplateStatusPublished),
		Version: maxInt(t.Version, 1), CreatedBy: t.CreatedBy, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func ptrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	if b == nil {
		return []byte("{}")
	}
	return b
}

func newOpenToken() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
