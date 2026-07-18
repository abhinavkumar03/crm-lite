package service

import (
	"context"
	"strings"

	"github.com/abhinavkumar03/crm-lite/backend/internal/search/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Search(
	ctx context.Context,
	orgID, query string,
) (*dto.SearchResponse, error) {
	query = strings.TrimSpace(query)
	hits, err := s.repo.Search(ctx, orgID, query, 15)
	if err != nil {
		return nil, err
	}
	if hits == nil {
		hits = []dto.SearchHit{}
	}
	return &dto.SearchResponse{Results: hits}, nil
}
