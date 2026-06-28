package service

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/entity"
	"github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
)

type Service struct {
	repository *repository.Repository
}

func New(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	ownerID string,
	req dto.CreateContactRequest,
) (*dto.ContactResponse, error) {

	contact := &entity.Contact{
		OwnerID:   ownerID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Company:   req.Company,
		JobTitle:  req.JobTitle,
		Notes:     req.Notes,
	}

	if err := s.repository.Create(ctx, contact); err != nil {
		return nil, err
	}

	return &dto.ContactResponse{
		ID:        contact.ID,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		Company:   contact.Company,
		JobTitle:  contact.JobTitle,
		Notes:     contact.Notes,
	}, nil
}

func (s *Service) List(
	ctx context.Context,
	ownerID string,
	req dto.ListContactsRequest,
) ([]dto.ContactResponse, error) {

	contacts, err := s.repository.List(
		ctx,
		ownerID,
		req,
	)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ContactResponse, 0, len(contacts))

	for _, contact := range contacts {
		response = append(response, dto.ContactResponse{
			ID:        contact.ID,
			FirstName: contact.FirstName,
			LastName:  contact.LastName,
			Email:     contact.Email,
			Phone:     contact.Phone,
			Company:   contact.Company,
			JobTitle:  contact.JobTitle,
			Notes:     contact.Notes,
		})
	}

	return response, nil
}
