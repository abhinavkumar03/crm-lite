package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/activity/entity"
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

func (r *Repository) Create(
	ctx context.Context,
	activity *entity.Activity,
) error {

	query := `
	INSERT INTO activities (

		id,

		entity_type,

		entity_id,

		action,

		description,

		performed_by,

		metadata,

		created_at

	)

	VALUES (

		$1,$2,$3,$4,$5,$6,$7,$8

	)
	`

	_, err := r.db.Exec(
		ctx,
		query,

		activity.ID,

		activity.EntityType,

		activity.EntityID,

		activity.Action,

		activity.Description,

		activity.PerformedBy,

		activity.Metadata,

		activity.CreatedAt,
	)

	return err
}

func (r *Repository) List(
	ctx context.Context,
	entityType string,
	entityID string,
) ([]entity.Activity, error) {

	query := `
	SELECT

		id,

		entity_type,

		entity_id,

		action,

		description,

		performed_by,

		metadata,

		created_at

	FROM activities

	WHERE entity_type=$1

	AND entity_id=$2

	ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		entityType,
		entityID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	activities := make([]entity.Activity, 0)

	for rows.Next() {

		var activity entity.Activity

		err := rows.Scan(
			&activity.ID,
			&activity.EntityType,
			&activity.EntityID,
			&activity.Action,
			&activity.Description,
			&activity.PerformedBy,
			&activity.Metadata,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		activities = append(
			activities,
			activity,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return activities, nil
}
