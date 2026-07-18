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
	ownerID string,
	refresh bool,
) (*dto.DashboardResponse, error) {
	key := cache.DashboardKey(ownerID)

	if !refresh {
		var cached dto.DashboardResponse
		if s.cache.GetJSON(ctx, key, &cached) {
			return &cached, nil
		}
	}

	data, err := s.repository.GetMetrics(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	data.RecentLeads, err = s.repository.RecentLeads(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	data.UpcomingTasks, err = s.repository.UpcomingTasks(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	data.RecentActivities, err = s.repository.RecentActivities(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	s.cache.SetJSON(ctx, key, data, cache.TTLMedium)
	return data, nil
}

// InvalidateDashboard drops the cached dashboard for a user so the next read
// rebuilds from Postgres.
func InvalidateDashboard(ctx context.Context, c *cache.Cache, userID string) {
	c.InvalidateDashboard(ctx, userID)
}
