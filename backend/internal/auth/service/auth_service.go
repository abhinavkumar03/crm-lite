package service

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/repository"
)

type AuthService struct {
	repository *repository.AuthRepository
}

func New(
	repository *repository.AuthRepository,
) *AuthService {

	return &AuthService{
		repository: repository,
	}
}
