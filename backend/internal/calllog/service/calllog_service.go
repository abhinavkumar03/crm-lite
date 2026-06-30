package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	activityEntity "github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/repository"
	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
	"github.com/google/uuid"
)

type Service struct {
	repository      *repository.Repository
	leadRepo        *leadRepository.Repository
	contactRepo     *contactRepository.Repository
	taskRepo        *taskRepository.Repository
	activityService *activityService.Service
}

func New(
	repository *repository.Repository,
	leadRepo *leadRepository.Repository,
	contactRepo *contactRepository.Repository,
	taskRepo *taskRepository.Repository,
	activityService *activityService.Service,

) *Service {

	return &Service{
		repository:      repository,
		leadRepo:        leadRepo,
		contactRepo:     contactRepo,
		taskRepo:        taskRepo,
		activityService: activityService,
	}
}

func (s *Service) validateOwnership(
	ctx context.Context,
	entityType string,
	entityID string,
	ownerID string,
) error {

	switch strings.ToUpper(entityType) {

	case "LEAD":

		ok, err := s.leadRepo.ExistsByIDAndOwner(
			ctx,
			entityID,
			ownerID,
		)

		if err != nil {
			return err
		}

		if !ok {
			return errors.New("lead not found")
		}

	case "CONTACT":

		ok, err := s.contactRepo.ExistsByIDAndOwner(
			ctx,
			entityID,
			ownerID,
		)

		if err != nil {
			return err
		}

		if !ok {
			return errors.New("contact not found")
		}

	case "TASK":

		ok, err := s.taskRepo.ExistsByIDAndOwner(
			ctx,
			entityID,
			ownerID,
		)

		if err != nil {
			return err
		}

		if !ok {
			return errors.New("task not found")
		}

	default:
		return errors.New("invalid entity type")
	}

	return nil
}

func validateEntityType(entityType string) error {

	switch strings.ToUpper(entityType) {

	case "LEAD",
		"CONTACT",
		"TASK":

		return nil

	default:

		return errors.New("invalid entity type")
	}
}

func validateDirection(direction string) error {

	switch strings.ToUpper(direction) {

	case "INCOMING",
		"OUTGOING":

		return nil

	default:

		return errors.New("invalid call direction")
	}
}

func validateStatus(status string) error {

	switch strings.ToUpper(status) {

	case "COMPLETED",
		"MISSED",
		"NO_ANSWER",
		"BUSY",
		"VOICEMAIL",
		"CANCELLED":

		return nil

	default:

		return errors.New("invalid call status")
	}
}

func (s *Service) Create(
	ctx context.Context,
	ownerID string,
	entityType string,
	entityID string,
	req dto.CreateCallLogRequest,
) error {

	if err := validateEntityType(entityType); err != nil {
		return err
	}

	if err := validateDirection(req.Direction); err != nil {
		return err
	}

	if err := validateStatus(req.Status); err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		entityType,
		entityID,
		ownerID,
	); err != nil {
		return err
	}

	now := time.Now().UTC()

	call := &entity.CallLog{
		ID:              uuid.NewString(),
		EntityType:      entity.EntityType(strings.ToUpper(entityType)),
		EntityID:        entityID,
		Direction:       entity.CallDirection(strings.ToUpper(req.Direction)),
		Status:          entity.CallStatus(strings.ToUpper(req.Status)),
		DurationSeconds: req.DurationSeconds,
		Summary:         strings.TrimSpace(req.Summary),
		FollowUpAt:      req.FollowUpAt,
		CreatedBy:       ownerID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repository.Create(ctx, call); err != nil {
		return err
	}

	metadata, _ := json.Marshal(map[string]any{
		"direction":        call.Direction,
		"status":           call.Status,
		"duration_seconds": call.DurationSeconds,
	})

	_ = s.activityService.Create(
		ctx,
		entityType,
		entityID,
		activityEntity.ActionCallLogged,
		"Logged a call",
		metadata,
		ownerID,
	)

	return nil
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	entityType string,
	entityID string,
) ([]dto.CallLogResponse, error) {

	if err := s.validateOwnership(
		ctx,
		entityType,
		entityID,
		ownerID,
	); err != nil {
		return nil, err
	}

	calls, err := s.repository.List(
		ctx,
		entityType,
		entityID,
	)

	if err != nil {
		return nil, err
	}

	response := make([]dto.CallLogResponse, 0, len(calls))

	for _, c := range calls {

		response = append(response, dto.CallLogResponse{
			ID:              c.ID,
			EntityType:      string(c.EntityType),
			EntityID:        c.EntityID,
			Direction:       string(c.Direction),
			Status:          string(c.Status),
			DurationSeconds: c.DurationSeconds,
			Summary:         c.Summary,
			FollowUpAt:      c.FollowUpAt,
			CreatedBy:       c.CreatedBy,
			UpdatedBy:       c.UpdatedBy,
			CreatedAt:       c.CreatedAt,
			UpdatedAt:       c.UpdatedAt,
		})
	}

	return response, nil
}

func (s *Service) Update(
	ctx context.Context,
	ownerID string,
	callLogID string,
	req dto.UpdateCallLogRequest,
) error {

	call, err := s.repository.GetByID(
		ctx,
		callLogID,
	)

	if err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		string(call.EntityType),
		call.EntityID,
		ownerID,
	); err != nil {
		return err
	}

	if req.Direction != "" {
		if err := validateDirection(req.Direction); err != nil {
			return err
		}
		call.Direction = entity.CallDirection(strings.ToUpper(req.Direction))
	}

	if req.Status != "" {
		if err := validateStatus(req.Status); err != nil {
			return err
		}
		call.Status = entity.CallStatus(strings.ToUpper(req.Status))
	}

	call.DurationSeconds = req.DurationSeconds
	call.Summary = strings.TrimSpace(req.Summary)
	call.FollowUpAt = req.FollowUpAt

	call.UpdatedAt = time.Now().UTC()
	call.UpdatedBy = &ownerID

	if err := s.repository.Update(ctx, call); err != nil {
		return err
	}

	_ = s.activityService.Create(
		ctx,
		string(call.EntityType),
		call.EntityID,
		activityEntity.ActionCallUpdated,
		"Updated a call log",
		nil,
		ownerID,
	)

	return nil
}

func (s *Service) Delete(
	ctx context.Context,
	ownerID string,
	callLogID string,
) error {

	call, err := s.repository.GetByID(
		ctx,
		callLogID,
	)

	if err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		string(call.EntityType),
		call.EntityID,
		ownerID,
	); err != nil {
		return err
	}

	if err := s.repository.Delete(
		ctx,
		callLogID,
	); err != nil {
		return err
	}

	_ = s.activityService.Create(
		ctx,
		string(call.EntityType),
		call.EntityID,
		activityEntity.ActionCallDeleted,
		"Deleted a call log",
		nil,
		ownerID,
	)

	return nil
}
