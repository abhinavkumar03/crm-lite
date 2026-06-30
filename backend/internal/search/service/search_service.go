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

		leadResults    = make([]leadDto.LeadResponse, 0)
		contactResults = make([]contactDto.ContactResponse, 0)
		taskResults    = make([]taskDto.TaskResponse, 0)

		leadErr    error
		contactErr error
		taskErr    error
	)

	wg.Add(3)

	go func() {
		defer wg.Done()

		leadResults, leadErr = s.leadRepo.Search(
			ctx,
			ownerID,
			query,
		)
	}()

	go func() {
		defer wg.Done()

		contactResults, contactErr = s.contactRepo.Search(
			ctx,
			ownerID,
			query,
		)
	}()

	go func() {
		defer wg.Done()

		taskResults, taskErr = s.taskRepo.Search(
			ctx,
			ownerID,
			query,
		)
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
