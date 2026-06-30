package service

import (
	"context"
	"errors"
	"strings"
	"time"

	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	noteRepository "github.com/abhinavkumar03/crm-lite/backend/internal/note/repository"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/note/dto"
	"github.com/google/uuid"
)

type Service struct {
	noteRepo    *noteRepository.Repository
	leadRepo    *leadRepository.Repository
	contactRepo *contactRepository.Repository
	taskRepo    *taskRepository.Repository
}

func New(
	noteRepo *noteRepository.Repository,
	leadRepo *leadRepository.Repository,
	contactRepo *contactRepository.Repository,
	taskRepo *taskRepository.Repository,
) *Service {

	return &Service{
		noteRepo:    noteRepo,
		leadRepo:    leadRepo,
		contactRepo: contactRepo,
		taskRepo:    taskRepo,
	}
}

func validateEntityType(entityType string) error {

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
	req dto.CreateNoteRequest,
) error {

	if err := validateEntityType(req.EntityType); err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		req.EntityType,
		req.EntityID,
		ownerID,
	); err != nil {
		return err
	}

	now := time.Now().UTC()

	note := &entity.Note{
		ID:         uuid.NewString(),
		EntityType: entity.EntityType(strings.ToUpper(req.EntityType)),
		EntityID:   req.EntityID,
		Note:       strings.TrimSpace(req.Note),
		CreatedBy:  ownerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return s.noteRepo.Create(
		ctx,
		note,
	)
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	entityType string,
	entityID string,
) ([]dto.NoteResponse, error) {

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

	notes, err := s.noteRepo.List(
		ctx,
		strings.ToUpper(entityType),
		entityID,
	)

	if err != nil {
		return nil, err
	}

	response := make([]dto.NoteResponse, 0, len(notes))

	for _, n := range notes {

		response = append(response, dto.NoteResponse{
			ID:         n.ID,
			EntityType: string(n.EntityType),
			EntityID:   n.EntityID,
			Note:       n.Note,
			CreatedBy:  n.CreatedBy,
			UpdatedBy:  n.UpdatedBy,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
		})
	}

	return response, nil
}

func (s *Service) Update(
	ctx context.Context,
	ownerID string,
	noteID string,
	req dto.UpdateNoteRequest,
) error {

	note, err := s.noteRepo.GetByID(
		ctx,
		noteID,
	)

	if err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		string(note.EntityType),
		note.EntityID,
		ownerID,
	); err != nil {
		return err
	}

	note.Note = strings.TrimSpace(req.Note)
	note.UpdatedAt = time.Now().UTC()
	note.UpdatedBy = &ownerID

	return s.noteRepo.Update(
		ctx,
		note,
	)
}

func (s *Service) Delete(
	ctx context.Context,
	ownerID string,
	noteID string,
) error {

	note, err := s.noteRepo.GetByID(
		ctx,
		noteID,
	)

	if err != nil {
		return err
	}

	if err := s.validateOwnership(
		ctx,
		string(note.EntityType),
		note.EntityID,
		ownerID,
	); err != nil {
		return err
	}

	return s.noteRepo.Delete(
		ctx,
		noteID,
	)
}
