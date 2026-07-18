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
	n.Status = entity.StatusQueued
	r.created = n
	return nil
}
func (r *fakeRepo) GetByID(_ context.Context, _, _ string) (*entity.Notification, error) {
	return r.created, nil
}
func (r *fakeRepo) List(_ context.Context, _ string, _ dto.ListQuery) ([]entity.Notification, int, error) {
	return nil, 0, nil
}
func (r *fakeRepo) MarkSent(_ context.Context, _, _ string) error { return nil }
func (r *fakeRepo) MarkFailed(_ context.Context, id, reason string) error {
	r.failedID = id
	r.failReason = reason
	return nil
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

	resp, err := svc.Send(context.Background(), "org1", "user1", dto.SendNotificationRequest{
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

	_, err := svc.Send(context.Background(), "org1", "user1", dto.SendNotificationRequest{
		Channel: "email",
		To:      "a@b.com",
		Body:    "hello",
	})
	if err == nil {
		t.Fatal("expected error when enqueue fails")
	}
	if repo.failedID != "n1" {
		t.Fatalf("expected notification to be marked failed, got %q", repo.failedID)
	}
}
