package service

import (
	"context"
	"encoding/json"

	activityEntity "github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	"github.com/abhinavkumar03/crm-lite/backend/internal/jobs"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
)

type Service struct {
	repository      *repository.Repository
	producer        *jobs.Producer
	activityService *activityService.Service
}

func New(
	repo *repository.Repository,
	producer *jobs.Producer,
	activityService *activityService.Service,

) *Service {

	return &Service{
		repository:      repo,
		producer:        producer,
		activityService: activityService,
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

	_ = s.activityService.Create(
		ctx,
		"LEAD",
		lead.ID,
		activityEntity.ActionLeadCreated,
		"Lead created",
		nil,
		ownerID,
	)

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
	req dto.ListLeadRequest,
) (*dto.ListLeadResponse, error) {

	return s.repository.List(
		ctx,
		ownerID,
		req,
	)
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

	oldStatus := lead.Status
	newStatus := req.Status

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

	_ = s.activityService.Create(
		ctx,
		"LEAD",
		lead.ID,
		activityEntity.ActionLeadUpdated,
		"Lead updated",
		nil,
		ownerID,
	)

	if oldStatus != newStatus {
		metadata, _ := json.Marshal(map[string]any{
			"from": oldStatus,
			"to":   newStatus,
		})

		_ = s.activityService.Create(
			ctx,
			"LEAD",
			lead.ID,
			activityEntity.ActionLeadStatusChanged,
			"Lead status changed",
			metadata,
			ownerID,
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

func (s *Service) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	deleted, err := s.repository.Delete(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return false, err
	}

	if deleted {
		_ = s.activityService.Create(
			ctx,
			"LEAD",
			id,
			activityEntity.ActionLeadDeleted,
			"Lead deleted",
			nil,
			ownerID,
		)
	}

	return deleted, nil
}
