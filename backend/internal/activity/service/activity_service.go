package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/repository"
	"github.com/google/uuid"

	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Service struct {
	repository *repository.Repository

	leadRepo *leadRepository.Repository

	contactRepo *contactRepository.Repository

	taskRepo *taskRepository.Repository
}

func New(
	repository *repository.Repository,
	leadRepo *leadRepository.Repository,
	contactRepo *contactRepository.Repository,
	taskRepo *taskRepository.Repository,
) *Service {

	return &Service{
		repository:  repository,
		leadRepo:    leadRepo,
		contactRepo: contactRepo,
		taskRepo:    taskRepo,
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
	}

	return nil
}

func (s *Service) Create(
	ctx context.Context,
	entityType string,
	entityID string,
	action string,
	description string,
	metadata []byte,
	performedBy string,
) error {

	if err := s.validateOwnership(
		ctx,
		entityType,
		entityID,
		performedBy,
	); err != nil {
		return err
	}

	activity := &entity.Activity{

		ID: uuid.NewString(),

		EntityType: entity.EntityType(entityType),

		EntityID: entityID,

		Action: action,

		Description: description,

		PerformedBy: performedBy,

		Metadata: metadata,

		CreatedAt: time.Now().UTC(),
	}

	return s.repository.Create(
		ctx,
		activity,
	)
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	entityType string,
	entityID string,
) ([]dto.ActivityResponse, error) {

	if err := s.validateOwnership(
		ctx,
		entityType,
		entityID,
		ownerID,
	); err != nil {
		return nil, err
	}

	activities, err := s.repository.List(
		ctx,
		entityType,
		entityID,
	)

	if err != nil {
		return nil, err
	}

	response := make(
		[]dto.ActivityResponse,
		0,
		len(activities),
	)

	for _, activity := range activities {

		response = append(
			response,
			dto.ActivityResponse{

				ID: activity.ID,

				Action: activity.Action,

				Description: activity.Description,

				PerformedBy: activity.PerformedBy,

				Metadata: activity.Metadata,

				CreatedAt: activity.CreatedAt,
			},
		)
	}

	return response, nil
}
