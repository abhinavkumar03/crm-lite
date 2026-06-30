package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/repository"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	repository *repository.Repository
	redis      *redis.Client
}

func New(
	repo *repository.Repository,
	redis *redis.Client,
) *Service {

	return &Service{
		repository: repo,
		redis:      redis,
	}
}

func (s *Service) GetDashboard(
	ctx context.Context,
	ownerID string,
	refresh bool,
) (*dto.DashboardResponse, error) {

	key := "dashboard:" + ownerID

	if !refresh {
		cached, err := s.redis.Get(ctx, key).Result()

		if err == nil {
			var response dto.DashboardResponse

			if json.Unmarshal([]byte(cached), &response) == nil {
				return &response, nil
			}
		}
	}

	data, err := s.repository.GetMetrics(
		ctx,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	data.RecentLeads, err = s.repository.RecentLeads(
		ctx,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	data.UpcomingTasks, err = s.repository.UpcomingTasks(
		ctx,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	bytes, _ := json.Marshal(data)

	_ = s.redis.Set(
		ctx,
		key,
		bytes,
		5*time.Minute,
	).Err()

	return data, nil
}

func InvalidateDashboard(
	ctx context.Context,
	redis *redis.Client,
	userID string,
) {

	key := "dashboard:" + userID

	redis.Del(
		ctx,
		key,
	)
}
