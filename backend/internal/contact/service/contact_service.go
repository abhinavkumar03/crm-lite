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
) (*dto.ListContactsResponse, error) {

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
) (*dto.ContactResponse, error) {

	contact, err := s.repository.GetByID(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, nil
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

func (s *Service) Update(
	ctx context.Context,
	id string,
	ownerID string,
	req dto.UpdateContactRequest,
) (*dto.ContactResponse, error) {

	contact, err := s.repository.GetByID(
		ctx,
		id,
		ownerID,
	)

	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, nil
	}

	if req.FirstName != "" {
		contact.FirstName = req.FirstName
	}

	if req.LastName != "" {
		contact.LastName = req.LastName
	}

	if req.Email != "" {
		contact.Email = req.Email
	}

	if req.Phone != "" {
		contact.Phone = req.Phone
	}

	if req.Company != "" {
		contact.Company = req.Company
	}

	if req.JobTitle != "" {
		contact.JobTitle = req.JobTitle
	}

	if req.Notes != "" {
		contact.Notes = req.Notes
	}

	err = s.repository.Update(
		ctx,
		contact,
	)

	if err != nil {
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
