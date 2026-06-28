package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/task/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	task *entity.Task,
) error {

	query := `
	INSERT INTO tasks
	(
		owner_id,
		lead_id,
		contact_id,
		title,
		description,
		status,
		due_date
	)
	VALUES
	(
		$1,$2,$3,$4,$5,$6,$7
	)
	RETURNING
		id,
		created_at,
		updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		task.OwnerID,
		task.LeadID,
		task.ContactID,
		task.Title,
		task.Description,
		task.Status,
		task.DueDate,
	).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
}
