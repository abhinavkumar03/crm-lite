package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	activityEntity "github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
	activityService "github.com/abhinavkumar03/crm-lite/backend/internal/activity/service"

	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"

	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/repository"
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

func validateEntityType(
	entityType string,
) error {

	switch strings.ToUpper(entityType) {

	case "LEAD":
		return nil

	case "CONTACT":
		return nil

	case "TASK":
		return nil

	default:
		return errors.New("invalid entity type")
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
	ownerID string,
	entityType string,
	entityID string,
	req dto.CreateAttachmentRequest,
) error {

	if err := validateEntityType(entityType); err != nil {
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

	attachment := &entity.Attachment{
		ID: uuid.NewString(),

		EntityType: entity.EntityType(
			strings.ToUpper(entityType),
		),

		EntityID: entityID,

		FileName: strings.TrimSpace(req.FileName),

		FileURL: req.FileURL,

		PublicID: req.PublicID,

		ResourceType: req.ResourceType,

		FileSize: req.FileSize,

		UploadedBy: ownerID,

		CreatedAt: time.Now().UTC(),
	}

	if err := s.repository.Create(
		ctx,
		attachment,
	); err != nil {
		return err
	}

	metadata, _ := json.Marshal(map[string]any{
		"file_name":     attachment.FileName,
		"resource_type": attachment.ResourceType,
	})

	_ = s.activityService.Create(
		ctx,
		entityType,
		entityID,
		activityEntity.ActionAttachmentAdded,
		"Uploaded attachment",
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
) ([]dto.AttachmentResponse, error) {

	if err := validateEntityType(entityType); err != nil {
		return nil, err
	}

	if err := s.validateOwnership(
		ctx,
		entityType,
		entityID,
		ownerID,
	); err != nil {
		return nil, err
	}

	attachments, err := s.repository.List(
		ctx,
		entityType,
		entityID,
	)

	if err != nil {
		return nil, err
	}

	response := make(
		[]dto.AttachmentResponse,
		0,
		len(attachments),
	)

	for _, attachment := range attachments {

		response = append(
			response,
			dto.AttachmentResponse{

				ID: attachment.ID,

				EntityType: string(
					attachment.EntityType,
				),

				EntityID: attachment.EntityID,

				FileName: attachment.FileName,

				FileURL: attachment.FileURL,

				PublicID: attachment.PublicID,

				ResourceType: attachment.ResourceType,

				FileSize: attachment.FileSize,

				UploadedBy: attachment.UploadedBy,

				CreatedAt: attachment.CreatedAt,
			},
		)
	}

	return response, nil
}

func (s *Service) Delete(
	ctx context.Context,
	ownerID string,
	attachmentID string,
) error {

	attachment, err := s.repository.GetByID(
		ctx,
		attachmentID,
	)

	if err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		string(attachment.EntityType),
		attachment.EntityID,
		ownerID,
	); err != nil {
		return err
	}

	// MVP:
	// We only remove the database record.
	// Cloudinary cleanup will be added later.

	if err := s.repository.Delete(
		ctx,
		attachmentID,
	); err != nil {
		return err
	}

	metadata, _ := json.Marshal(map[string]any{
		"file_name": attachment.FileName,
	})

	_ = s.activityService.Create(
		ctx,
		string(attachment.EntityType),
		attachment.EntityID,
		activityEntity.ActionAttachmentDeleted,
		"Deleted attachment",
		metadata,
		ownerID,
	)

	return nil
}
