package service

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/dashboard/repository"
)

type Service struct {
	repository *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) GetDashboard(
	ctx context.Context,
	ownerID string,
) (*dto.DashboardResponse, error) {

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

	return data, nil
}
