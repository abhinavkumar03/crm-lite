package service

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
)

type Service struct {
	repository *repository.Repository
}

func New(
	repository *repository.Repository,
) *Service {

	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	ownerID string,
	req dto.CreateLeadRequest,
) (*dto.LeadResponse, error) {

	lead := &entity.Lead{
		OwnerID: ownerID,
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Company: req.Company,
		Status:  entity.StatusNew,
		Notes:   req.Notes,
	}

	err := s.repository.Create(
		ctx,
		lead,
	)

	if err != nil {
		return nil, err
	}

	return &dto.LeadResponse{
		ID:      lead.ID,
		Name:    lead.Name,
		Email:   lead.Email,
		Phone:   lead.Phone,
		Company: lead.Company,
		Status:  lead.Status,
		Notes:   lead.Notes,
	}, nil
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	req dto.ListLeadsRequest,
) ([]dto.LeadResponse, error) {

	leads, err := s.repository.List(
		ctx,
		ownerID,
		req,
	)

	if err != nil {
		return nil, err
	}

	response := make([]dto.LeadResponse, 0, len(leads))

	for _, lead := range leads {

		response = append(response, dto.LeadResponse{
			ID:      lead.ID,
			Name:    lead.Name,
			Email:   lead.Email,
			Phone:   lead.Phone,
			Company: lead.Company,
			Status:  lead.Status,
			Notes:   lead.Notes,
		})
	}

	return response, nil
}
