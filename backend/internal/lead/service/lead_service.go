package service

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
)

type Service struct {
	repository *repository.Repository
	producer   *jobs.Producer
}

func New(
	repo *repository.Repository,
	producer *jobs.Producer,
) *Service {

	return &Service{
		repository: repo,
		producer:   producer,
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

	if err := s.producer.Publish(
		ctx,
		jobs.Job{
			Type:   jobs.JobLeadCreated,
			UserID: ownerID,
			Payload: map[string]interface{}{
				"lead_id": lead.ID,
				"name":    lead.Name,
			},
		},
	); err != nil {
		// log if desired, but don't fail the HTTP request
	}

	if lead.Email != "" {
		_ = s.producer.Publish(
			ctx,
			jobs.Job{
				Type:   jobs.JobSendEmail,
				UserID: ownerID,
				Payload: map[string]interface{}{
					"email": lead.Email,
					"name":  lead.Name,
				},
			},
		)
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

func (s *Service) GetByID(
	ctx context.Context,
	id string,
	ownerID string,
) (*dto.LeadResponse, error) {

	lead, err := s.repository.GetByID(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	if lead == nil {
		return nil, nil
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

func (s *Service) Update(
	ctx context.Context,
	id string,
	ownerID string,
	req dto.UpdateLeadRequest,
) (*dto.LeadResponse, error) {

	lead, err := s.repository.GetByID(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	if lead == nil {
		return nil, nil
	}

	if req.Name != "" {
		lead.Name = req.Name
	}

	if req.Email != "" {
		lead.Email = req.Email
	}

	if req.Phone != "" {
		lead.Phone = req.Phone
	}

	if req.Company != "" {
		lead.Company = req.Company
	}

	if req.Status != "" {
		lead.Status = req.Status
	}

	if req.Notes != "" {
		lead.Notes = req.Notes
	}

	err = s.repository.Update(
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

func (s *Service) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	return s.repository.Delete(
		ctx,
		id,
		ownerID,
	)
}
