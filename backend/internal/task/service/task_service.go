package service

import (
	"context"
	"errors"
	"time"

	contactrepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadrepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Service struct {
	taskRepository    *repository.Repository
	leadRepository    *leadrepository.Repository
	contactRepository *contactrepository.Repository
}

func New(
	taskRepo *repository.Repository,
	leadRepo *leadrepository.Repository,
	contactRepo *contactrepository.Repository,
) *Service {

	return &Service{
		taskRepository:    taskRepo,
		leadRepository:    leadRepo,
		contactRepository: contactRepo,
	}
}

func (s *Service) Create(
	ctx context.Context,
	ownerID string,
	req dto.CreateTaskRequest,
) (*dto.TaskResponse, error) {

	// Validate Lead ownership (optional)
	if req.LeadID != nil {

		lead, err := s.leadRepository.GetByID(
			ctx,
			*req.LeadID,
			ownerID,
		)

		if err != nil {
			return nil, err
		}

		if lead == nil {
			return nil, errors.New("lead not found")
		}
	}

	// Validate Contact ownership (optional)
	if req.ContactID != nil {

		contact, err := s.contactRepository.GetByID(
			ctx,
			*req.ContactID,
			ownerID,
		)

		if err != nil {
			return nil, err
		}

		if contact == nil {
			return nil, errors.New("contact not found")
		}
	}

	task := &entity.Task{
		OwnerID:     ownerID,
		LeadID:      req.LeadID,
		ContactID:   req.ContactID,
		Title:       req.Title,
		Description: req.Description,
		Status:      entity.StatusPending,
	}

	if req.DueDate != nil {
		t, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date")
		}
		task.DueDate = &t
	}

	if err := s.taskRepository.Create(ctx, task); err != nil {
		return nil, err
	}

	var dueDate *string
	if task.DueDate != nil {
		v := task.DueDate.Format(time.RFC3339)
		dueDate = &v
	}

	return &dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		LeadID:      task.LeadID,
		ContactID:   task.ContactID,
		DueDate:     dueDate,
	}, nil
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	req dto.ListTasksRequest,
) ([]dto.TaskResponse, error) {

	tasks, err := s.taskRepository.List(
		ctx,
		ownerID,
		req,
	)

	if err != nil {
		return nil, err
	}

	response := make([]dto.TaskResponse, 0, len(tasks))

	for _, task := range tasks {

		var dueDate *string

		if task.DueDate != nil {
			v := task.DueDate.Format(time.RFC3339)
			dueDate = &v
		}

		response = append(response, dto.TaskResponse{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			LeadID:      task.LeadID,
			ContactID:   task.ContactID,
			DueDate:     dueDate,
		})
	}

	return response, nil
}

func (s *Service) GetByID(
	ctx context.Context,
	id string,
	ownerID string,
) (*dto.TaskResponse, error) {

	task, err := s.taskRepository.GetByID(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, nil
	}

	var dueDate *string

	if task.DueDate != nil {
		value := task.DueDate.Format(time.RFC3339)
		dueDate = &value
	}

	return &dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		LeadID:      task.LeadID,
		ContactID:   task.ContactID,
		DueDate:     dueDate,
	}, nil
}

func (s *Service) Update(
	ctx context.Context,
	id string,
	ownerID string,
	req dto.UpdateTaskRequest,
) (*dto.TaskResponse, error) {

	task, err := s.taskRepository.GetByID(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, nil
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != "" {
		task.Description = req.Description
	}

	if req.Status != "" {
		task.Status = req.Status
	}

	if req.LeadID != nil {

		lead, err := s.leadRepository.GetByID(
			ctx,
			*req.LeadID,
			ownerID,
		)

		if err != nil {
			return nil, err
		}

		if lead == nil {
			return nil, errors.New("lead not found")
		}

		task.LeadID = req.LeadID
	}

	if req.ContactID != nil {

		contact, err := s.contactRepository.GetByID(
			ctx,
			*req.ContactID,
			ownerID,
		)

		if err != nil {
			return nil, err
		}

		if contact == nil {
			return nil, errors.New("contact not found")
		}

		task.ContactID = req.ContactID
	}

	if req.DueDate != nil {

		t, err := time.Parse(
			time.RFC3339,
			*req.DueDate,
		)

		if err != nil {
			return nil, errors.New("invalid due_date")
		}

		task.DueDate = &t
	}

	err = s.taskRepository.Update(
		ctx,
		task,
	)

	if err != nil {
		return nil, err
	}

	var dueDate *string

	if task.DueDate != nil {
		v := task.DueDate.Format(time.RFC3339)
		dueDate = &v
	}

	return &dto.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		LeadID:      task.LeadID,
		ContactID:   task.ContactID,
		DueDate:     dueDate,
	}, nil
}

func (s *Service) Delete(
	ctx context.Context,
	id string,
	ownerID string,
) (bool, error) {

	return s.taskRepository.Delete(
		ctx,
		id,
		ownerID,
	)
}
