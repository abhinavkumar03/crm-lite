package service

import (
	"context"
	"sync"

	contactDto "github.com/abhinavkumar03/crm-lite/backend/internal/contact/dto"
	contactRepository "github.com/abhinavkumar03/crm-lite/backend/internal/contact/repository"
	leadDto "github.com/abhinavkumar03/crm-lite/backend/internal/lead/dto"
	leadRepository "github.com/abhinavkumar03/crm-lite/backend/internal/lead/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/search/dto"
	taskDto "github.com/abhinavkumar03/crm-lite/backend/internal/task/dto"
	taskRepository "github.com/abhinavkumar03/crm-lite/backend/internal/task/repository"
)

type Service struct {
	leadRepo    *leadRepository.Repository
	contactRepo *contactRepository.Repository
	taskRepo    *taskRepository.Repository
}

func New(
	leadRepo *leadRepository.Repository,
	contactRepo *contactRepository.Repository,
	taskRepo *taskRepository.Repository,
) *Service {

	return &Service{
		leadRepo:    leadRepo,
		contactRepo: contactRepo,
		taskRepo:    taskRepo,
	}
}

func (s *Service) Search(
	ctx context.Context,
	ownerID string,
	query string,
) (*dto.SearchResponse, error) {

	var (
		wg sync.WaitGroup

		leadResults    []leadDto.LeadResponse
		contactResults []contactDto.ContactResponse
		taskResults    []taskDto.TaskResponse

		leadErr    error
		contactErr error
		taskErr    error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()

		res, err := s.leadRepo.List(
			ctx,
			ownerID,
			leadDto.ListLeadRequest{
				Page:   1,
				Limit:  5,
				Search: query,
			},
		)

		if err != nil {
			leadErr = err
			return
		}

		leadResults = res.Data
	}()

	go func() {
		defer wg.Done()

		res, err := s.contactRepo.List(
			ctx,
			ownerID,
			contactDto.ListContactsRequest{
				Page:   1,
				Limit:  5,
				Search: query,
			},
		)

		if err != nil {
			contactErr = err
			return
		}

		contactResults = res.Data
	}()

	go func() {
		defer wg.Done()

		res, err := s.taskRepo.List(
			ctx,
			ownerID,
			taskDto.ListTasksRequest{
				Page:   1,
				Limit:  5,
				Search: query,
			},
		)

		if err != nil {
			taskErr = err
			return
		}

		taskResults = res.Data
	}()

	wg.Wait()

	if leadErr != nil {
		return nil, leadErr
	}

	if contactErr != nil {
		return nil, contactErr
	}

	if taskErr != nil {
		return nil, taskErr
	}

	return &dto.SearchResponse{
		Leads:    leadResults,
		Contacts: contactResults,
		Tasks:    taskResults,
	}, nil
}
