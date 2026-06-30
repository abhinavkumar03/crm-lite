package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/search/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SearchLeads(
	ctx context.Context,
	ownerID string,
	query string,
) ([]dto.LeadResult, error)

func (r *Repository) SearchContacts(
	ctx context.Context,
	ownerID string,
	query string,
) ([]dto.ContactResult, error)

func (r *Repository) SearchTasks(
	ctx context.Context,
	ownerID string,
	query string,
) ([]dto.TaskResult, error)
