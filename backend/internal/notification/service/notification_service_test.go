package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hibiken/asynq"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
)

type fakeRepo struct {
	created    *entity.Notification
	failedID   string
	failReason string
}

func (r *fakeRepo) Create(_ context.Context, n *entity.Notification) error {
	n.ID = "n1"
	r.created = n
	return nil
}
func (r *fakeRepo) UpdateDraft(_ context.Context, _ *entity.Notification) error { return nil }
func (r *fakeRepo) GetByID(_ context.Context, _, _ string) (*entity.Notification, error) {
	return r.created, nil
}
func (r *fakeRepo) List(_ context.Context, _ string, _ dto.ListQuery) ([]entity.Notification, int, error) {
	return nil, 0, nil
}
func (r *fakeRepo) MarkSent(_ context.Context, _, _, _ string, _ map[string]any) error { return nil }
func (r *fakeRepo) MarkFailed(_ context.Context, id, reason string) error {
	r.failedID = id
	r.failReason = reason
	return nil
}
func (r *fakeRepo) MarkRetrying(_ context.Context, _, _ string) error { return nil }
func (r *fakeRepo) MarkQueued(_ context.Context, _ string) error       { return nil }
func (r *fakeRepo) CancelScheduled(_ context.Context, _, _ string) (*entity.Notification, error) {
	return nil, nil
}
func (r *fakeRepo) PromoteScheduled(_ context.Context, _ string) error { return nil }
func (r *fakeRepo) ListDueScheduled(_ context.Context, _ int) ([]entity.Notification, error) {
	return nil, nil
}
func (r *fakeRepo) AddDeliveryEvent(_ context.Context, _ *entity.DeliveryEvent) error { return nil }
func (r *fakeRepo) ListDeliveryEvents(_ context.Context, _, _ string) ([]entity.DeliveryEvent, error) {
	return nil, nil
}
func (r *fakeRepo) Metrics(_ context.Context, _ string) (*dto.MetricsResponse, error) {
	return &dto.MetricsResponse{}, nil
}
func (r *fakeRepo) LoadMergeContext(_ context.Context, _, _, _, _ string) (map[string]any, error) {
	return map[string]any{}, nil
}
func (r *fakeRepo) LinkAttachments(_ context.Context, _ string, _ []string) error { return nil }
func (r *fakeRepo) CreateTemplate(_ context.Context, _ *entity.Template) error { return nil }
func (r *fakeRepo) GetTemplate(_ context.Context, _, _ string) (*entity.Template, error) {
	return nil, nil
}
func (r *fakeRepo) UpdateTemplate(_ context.Context, _ *entity.Template) error { return nil }
func (r *fakeRepo) PublishTemplate(_ context.Context, _, _ string) (*entity.Template, error) {
	return nil, nil
}
func (r *fakeRepo) DeleteTemplate(_ context.Context, _, _ string) error         { return nil }
func (r *fakeRepo) ListTemplates(_ context.Context, _, _, _ string, _, _ int) ([]entity.Template, int, error) {
	return nil, 0, nil
}

type fakeEnqueuer struct {
	published *jobs.Job
	err       error
}

func (e *fakeEnqueuer) Publish(_ context.Context, job jobs.Job, _ ...asynq.Option) error {
	if e.err != nil {
		return e.err
	}
	e.published = &job
	return nil
}

func TestSend_RendersPersistsAndEnqueues(t *testing.T) {
	repo := &fakeRepo{}
	enq := &fakeEnqueuer{}
	svc := New(repo, enq)

	resp, err := svc.Send(context.Background(), "org1", "user1", dto.ComposeRequest{
		Channel: "whatsapp",
		To:      "+15551234567",
		Body:    "Hi {{name}}, your quote is {{amount}}",
		Data:    map[string]any{"name": "Dana", "amount": 4200},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body != "Hi Dana, your quote is 4200" {
		t.Fatalf("body not rendered: %q", resp.Body)
	}
	if resp.Status != entity.StatusQueued {
		t.Fatalf("expected queued status, got %q", resp.Status)
	}
	if repo.created == nil || repo.created.CreatedBy == nil || *repo.created.CreatedBy != "user1" {
		t.Fatalf("created_by should be acting user")
	}
	if enq.published == nil || enq.published.Type != jobs.JobSendNotification {
		t.Fatalf("expected notification.send job, got %#v", enq.published)
	}
	if enq.published.Payload["notification_id"] != "n1" || enq.published.Payload["org_id"] != "org1" {
		t.Fatalf("unexpected job payload: %#v", enq.published.Payload)
	}
}

func TestSend_MarksFailedWhenEnqueueFails(t *testing.T) {
	repo := &fakeRepo{}
	enq := &fakeEnqueuer{err: errors.New("redis down")}
	svc := New(repo, enq)

	_, err := svc.Send(context.Background(), "org1", "user1", dto.ComposeRequest{
		Channel: "email",
		To:      "a@b.com",
		Subject: "Hello",
		Body:    "hello",
	})
	if err == nil {
		t.Fatal("expected error when enqueue fails")
	}
	if repo.failedID != "n1" {
		t.Fatalf("expected notification to be marked failed, got %q", repo.failedID)
	}
}
