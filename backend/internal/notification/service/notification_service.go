package service

import (
	"context"
	"encoding/json"
	"math"

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

// Repository is the persistence contract for notifications.
type Repository interface {
	Create(ctx context.Context, n *entity.Notification) error
	GetByID(ctx context.Context, orgID, id string) (*entity.Notification, error)
	List(ctx context.Context, orgID string, q dto.ListQuery) ([]entity.Notification, int, error)
	MarkSent(ctx context.Context, id, provider string) error
	MarkFailed(ctx context.Context, id, reason string) error
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

// Send renders, persists (status=queued) and enqueues a notification. Delivery
// happens asynchronously in the worker so the request returns immediately.
func (s *Service) Send(ctx context.Context, orgID, userID string, req dto.SendNotificationRequest) (*dto.NotificationResponse, error) {
	data := req.Data
	if data == nil {
		data = map[string]any{}
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	n := &entity.Notification{
		OrganizationID: orgID,
		Channel:        req.Channel,
		Recipient:      req.To,
		Subject:        ptrOrNil(notify.Render(req.Subject, data)),
		Body:           notify.Render(req.Body, data),
		Template:       ptrOrNil(req.Template),
		Data:           dataBytes,
		Status:         entity.StatusQueued,
		EntityType:     ptrOrNil(req.EntityType),
		EntityID:       ptrOrNil(req.EntityID),
		CreatedBy:      &userID,
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}

	job := jobs.Job{
		Type:   jobs.JobSendNotification,
		UserID: userID,
		Payload: map[string]interface{}{
			"notification_id": n.ID,
			"org_id":          orgID,
		},
	}
	if err := s.enqueuer.Publish(ctx, job); err != nil {
		_ = s.repo.MarkFailed(ctx, n.ID, "failed to enqueue: "+err.Error())
		return nil, err
	}

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

func toResponse(n *entity.Notification) dto.NotificationResponse {
	data := map[string]any{}
	if len(n.Data) > 0 {
		_ = json.Unmarshal(n.Data, &data)
	}
	return dto.NotificationResponse{
		ID:         n.ID,
		Channel:    n.Channel,
		Recipient:  n.Recipient,
		Subject:    n.Subject,
		Body:       n.Body,
		Template:   n.Template,
		Data:       data,
		Status:     n.Status,
		Provider:   n.Provider,
		Error:      n.Error,
		EntityType: n.EntityType,
		EntityID:   n.EntityID,
		CreatedBy:  n.CreatedBy,
		SentAt:     n.SentAt,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
}

func ptrOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
