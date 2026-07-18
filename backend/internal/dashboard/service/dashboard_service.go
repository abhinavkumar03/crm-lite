package service

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/shared/cache"
)

type Service struct {
	repository *repository.Repository
	cache      *cache.Cache
}

func New(repo *repository.Repository, c *cache.Cache) *Service {
	return &Service{repository: repo, cache: c}
}

func (s *Service) GetDashboard(
	ctx context.Context,
	orgID string,
	refresh bool,
) (*dto.DashboardResponse, error) {
	key := cache.DashboardKey(orgID)

	if !refresh {
		var cached dto.DashboardResponse
		if s.cache.GetJSON(ctx, key, &cached) {
			return &cached, nil
		}
	}

	data, err := s.repository.GetDashboard(ctx, orgID)
	if err != nil {
		return nil, err
	}

	s.cache.SetJSON(ctx, key, data, cache.TTLMedium)
	return data, nil
}

// InvalidateDashboard drops the cached dashboard for an organization.
func InvalidateDashboard(ctx context.Context, c *cache.Cache, orgID string) {
	c.InvalidateDashboard(ctx, orgID)
}
