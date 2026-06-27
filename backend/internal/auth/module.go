package auth

import (
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/repository"
	"github.com/abhinavkumar03/crm-lite/backend/internal/auth/service"
)

type Module struct {
	Service service.Service
	Repo    repository.Repository
}
