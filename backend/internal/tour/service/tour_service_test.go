package service

import (
	"context"
	"testing"

	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/entity"
)

// fakeRepo is an in-memory Repository keyed by (org, user, tourKey).
type fakeRepo struct {
	store map[string]*entity.TourProgress
}

func newFakeRepo() *fakeRepo { return &fakeRepo{store: map[string]*entity.TourProgress{}} }

func key(org, user, tourKey string) string { return org + "|" + user + "|" + tourKey }

func (f *fakeRepo) GetByUser(_ context.Context, org, user, tourKey string) (*entity.TourProgress, error) {
	return f.store[key(org, user, tourKey)], nil
}

func (f *fakeRepo) Upsert(_ context.Context, p *entity.TourProgress) error {
	clone := *p
	f.store[key(p.OrganizationID, p.UserID, p.TourKey)] = &clone
	return nil
}

func (f *fakeRepo) Restart(_ context.Context, org, user, tourKey string) (*entity.TourProgress, error) {
	p := &entity.TourProgress{
		OrganizationID: org,
		UserID:         user,
		TourKey:        tourKey,
		Status:         entity.StatusActive,
		CurrentStep:    0,
		CompletedSteps: []string{},
	}
	f.store[key(org, user, tourKey)] = p
	return p, nil
}

func TestGetDefaultsToActiveForNewUser(t *testing.T) {
	svc := New(newFakeRepo())

	resp, err := svc.Get(context.Background(), "org1", "user1", "")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if resp.TourKey != entity.DefaultTourKey {
		t.Errorf("expected default tour key %q, got %q", entity.DefaultTourKey, resp.TourKey)
	}
	if resp.Status != entity.StatusActive {
		t.Errorf("expected status active, got %q", resp.Status)
	}
	if resp.CurrentStep != 0 {
		t.Errorf("expected current step 0, got %d", resp.CurrentStep)
	}
}

func TestSaveAdvancesAndPersists(t *testing.T) {
	repo := newFakeRepo()
	svc := New(repo)

	step := 2
	resp, err := svc.Save(context.Background(), "org1", "user1", dto.UpdateProgressRequest{
		CurrentStep:    &step,
		CompletedSteps: []string{"welcome", "sidebar"},
	})
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if resp.CurrentStep != 2 {
		t.Errorf("expected current step 2, got %d", resp.CurrentStep)
	}
	if resp.Status != entity.StatusActive {
		t.Errorf("expected status active, got %q", resp.Status)
	}
	if resp.CompletedAt != nil {
		t.Errorf("expected nil completed_at while active")
	}

	stored, _ := repo.GetByUser(context.Background(), "org1", "user1", entity.DefaultTourKey)
	if stored == nil || stored.CurrentStep != 2 {
		t.Fatalf("expected persisted step 2, got %+v", stored)
	}
}

func TestSaveCompletedStampsCompletedAt(t *testing.T) {
	svc := New(newFakeRepo())

	resp, err := svc.Save(context.Background(), "org1", "user1", dto.UpdateProgressRequest{
		Status: entity.StatusCompleted,
	})
	if err != nil {
		t.Fatalf("Save returned error: %v", err)
	}
	if resp.Status != entity.StatusCompleted {
		t.Errorf("expected status completed, got %q", resp.Status)
	}
	if resp.CompletedAt == nil {
		t.Errorf("expected completed_at to be stamped on completion")
	}
}

func TestSaveReactivateClearsCompletedAt(t *testing.T) {
	svc := New(newFakeRepo())
	ctx := context.Background()

	if _, err := svc.Save(ctx, "org1", "user1", dto.UpdateProgressRequest{Status: entity.StatusCompleted}); err != nil {
		t.Fatalf("Save(completed) error: %v", err)
	}

	resp, err := svc.Save(ctx, "org1", "user1", dto.UpdateProgressRequest{Status: entity.StatusActive})
	if err != nil {
		t.Fatalf("Save(active) error: %v", err)
	}
	if resp.CompletedAt != nil {
		t.Errorf("expected completed_at cleared when reactivated")
	}
}

func TestRestartResetsProgress(t *testing.T) {
	svc := New(newFakeRepo())
	ctx := context.Background()

	step := 5
	if _, err := svc.Save(ctx, "org1", "user1", dto.UpdateProgressRequest{CurrentStep: &step, Status: entity.StatusCompleted}); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	resp, err := svc.Restart(ctx, "org1", "user1", dto.RestartRequest{})
	if err != nil {
		t.Fatalf("Restart error: %v", err)
	}
	if resp.CurrentStep != 0 || resp.Status != entity.StatusActive {
		t.Errorf("expected reset to step 0/active, got step=%d status=%q", resp.CurrentStep, resp.Status)
	}
}
