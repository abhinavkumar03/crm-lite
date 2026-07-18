package service

import (
	"context"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/tour/entity"
)

// Repository is the persistence contract for tour progress.
type Repository interface {
	GetByUser(ctx context.Context, orgID, userID, tourKey string) (*entity.TourProgress, error)
	Upsert(ctx context.Context, p *entity.TourProgress) error
	Restart(ctx context.Context, orgID, userID, tourKey string) (*entity.TourProgress, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// Get returns the user's progress for a tour. A user who has never started the
// tour gets a synthesized "active" default (not persisted) so the client can
// begin without a prior write.
func (s *Service) Get(ctx context.Context, orgID, userID, tourKey string) (*dto.ProgressResponse, error) {
	tourKey = normalizeKey(tourKey)

	p, err := s.repo.GetByUser(ctx, orgID, userID, tourKey)
	if err != nil {
		return nil, err
	}
	if p == nil {
		resp := defaultProgress(tourKey)
		return &resp, nil
	}

	resp := toResponse(p)
	return &resp, nil
}

// Save upserts the user's progress. It is the single write path for advancing,
// completing and skipping. completed_at is stamped when the tour first reaches a
// terminal (completed) state and cleared when it returns to active.
func (s *Service) Save(ctx context.Context, orgID, userID string, req dto.UpdateProgressRequest) (*dto.ProgressResponse, error) {
	tourKey := normalizeKey(req.TourKey)

	existing, err := s.repo.GetByUser(ctx, orgID, userID, tourKey)
	if err != nil {
		return nil, err
	}

	p := &entity.TourProgress{
		OrganizationID: orgID,
		UserID:         userID,
		TourKey:        tourKey,
		Status:         entity.StatusActive,
	}
	if existing != nil {
		p.Status = existing.Status
		p.CurrentStep = existing.CurrentStep
		p.CompletedSteps = existing.CompletedSteps
		p.StartedAt = existing.StartedAt
		p.CompletedAt = existing.CompletedAt
	}

	if req.Status != "" {
		p.Status = req.Status
	}
	if req.CurrentStep != nil {
		p.CurrentStep = *req.CurrentStep
	}
	if req.CompletedSteps != nil {
		p.CompletedSteps = req.CompletedSteps
	}

	// Keep completed_at consistent with the lifecycle.
	if p.Status == entity.StatusCompleted {
		if p.CompletedAt == nil {
			now := time.Now()
			p.CompletedAt = &now
		}
	} else {
		p.CompletedAt = nil
	}

	if err := s.repo.Upsert(ctx, p); err != nil {
		return nil, err
	}

	resp := toResponse(p)
	return &resp, nil
}

// Restart resets the tour to its first step.
func (s *Service) Restart(ctx context.Context, orgID, userID string, req dto.RestartRequest) (*dto.ProgressResponse, error) {
	p, err := s.repo.Restart(ctx, orgID, userID, normalizeKey(req.TourKey))
	if err != nil {
		return nil, err
	}
	resp := toResponse(p)
	return &resp, nil
}

func normalizeKey(key string) string {
	if key == "" {
		return entity.DefaultTourKey
	}
	return key
}

func defaultProgress(tourKey string) dto.ProgressResponse {
	now := time.Now()
	return dto.ProgressResponse{
		TourKey:        tourKey,
		Status:         entity.StatusActive,
		CurrentStep:    0,
		CompletedSteps: []string{},
		StartedAt:      now,
		CompletedAt:    nil,
		UpdatedAt:      now,
	}
}

func toResponse(p *entity.TourProgress) dto.ProgressResponse {
	steps := p.CompletedSteps
	if steps == nil {
		steps = []string{}
	}
	return dto.ProgressResponse{
		TourKey:        p.TourKey,
		Status:         p.Status,
		CurrentStep:    p.CurrentStep,
		CompletedSteps: steps,
		StartedAt:      p.StartedAt,
		CompletedAt:    p.CompletedAt,
		UpdatedAt:      p.UpdatedAt,
	}
}
