package repository

import (
	"context"
	"time"

	"github.com/abhinavkumar03/crm-lite/backend/internal/calllog/entity"
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
	call *entity.CallLog,
) error {

	query := `
INSERT INTO call_logs (

	id,

	entity_type,

	entity_id,

	direction,

	status,

	duration_seconds,

	summary,

	follow_up_at,

	created_by,

	updated_by,

	created_at,

	updated_at

)

VALUES (

	$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12

)
`

	_, err := r.db.Exec(

		ctx,

		query,

		call.ID,

		call.EntityType,

		call.EntityID,

		call.Direction,

		call.Status,

		call.DurationSeconds,

		call.Summary,

		call.FollowUpAt,

		call.CreatedBy,

		call.UpdatedBy,

		call.CreatedAt,

		call.UpdatedAt,
	)

	return err
}

func (r *Repository) List(
	ctx context.Context,
	entityType string,
	entityID string,
) ([]entity.CallLog, error) {

	query := `
SELECT

	id,

	entity_type,

	entity_id,

	direction,

	status,

	duration_seconds,

	summary,

	follow_up_at,

	created_by,

	updated_by,

	created_at,

	updated_at

FROM call_logs

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

	calls := make([]entity.CallLog, 0)

	for rows.Next() {

		var call entity.CallLog

		err := rows.Scan(

			&call.ID,

			&call.EntityType,

			&call.EntityID,

			&call.Direction,

			&call.Status,

			&call.DurationSeconds,

			&call.Summary,

			&call.FollowUpAt,

			&call.CreatedBy,

			&call.UpdatedBy,

			&call.CreatedAt,

			&call.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		calls = append(
			calls,
			call,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return calls, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*entity.CallLog, error) {

	query := `
SELECT

	id,

	entity_type,

	entity_id,

	direction,

	status,

	duration_seconds,

	summary,

	follow_up_at,

	created_by,

	updated_by,

	created_at,

	updated_at

FROM call_logs

WHERE id=$1
`

	var call entity.CallLog

	err := r.db.QueryRow(

		ctx,

		query,

		id,
	).Scan(

		&call.ID,

		&call.EntityType,

		&call.EntityID,

		&call.Direction,

		&call.Status,

		&call.DurationSeconds,

		&call.Summary,

		&call.FollowUpAt,

		&call.CreatedBy,

		&call.UpdatedBy,

		&call.CreatedAt,

		&call.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &call, nil
}

func (r *Repository) Update(
	ctx context.Context,
	call *entity.CallLog,
) error {

	query := `
UPDATE call_logs

SET

	direction=$1,

	status=$2,

	duration_seconds=$3,

	summary=$4,

	follow_up_at=$5,

	updated_by=$6,

	updated_at=$7

WHERE id=$8
`

	_, err := r.db.Exec(

		ctx,

		query,

		call.Direction,

		call.Status,

		call.DurationSeconds,

		call.Summary,

		call.FollowUpAt,

		call.UpdatedBy,

		call.UpdatedAt,

		call.ID,
	)

	return err
}

func (r *Repository) Delete(
	ctx context.Context,
	id string,
) error {

	query := `
DELETE

FROM call_logs

WHERE id=$1
`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}

func (r *Repository) GetLatestCall(
	ctx context.Context,
	entityType string,
	entityID string,
) (*entity.CallLog, error) {

	query := `
SELECT

	id,

	entity_type,

	entity_id,

	direction,

	status,

	duration_seconds,

	summary,

	follow_up_at,

	created_by,

	updated_by,

	created_at,

	updated_at

FROM call_logs

WHERE entity_type=$1

AND entity_id=$2

ORDER BY created_at DESC

LIMIT 1
`

	var call entity.CallLog

	err := r.db.QueryRow(
		ctx,
		query,
		entityType,
		entityID,
	).Scan(
		&call.ID,
		&call.EntityType,
		&call.EntityID,
		&call.Direction,
		&call.Status,
		&call.DurationSeconds,
		&call.Summary,
		&call.FollowUpAt,
		&call.CreatedBy,
		&call.UpdatedBy,
		&call.CreatedAt,
		&call.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &call, nil
}

func (r *Repository) GetUpcomingFollowUps(
	ctx context.Context,
	ownerID string,
	from time.Time,
	to time.Time,
) ([]entity.CallLog, error) {

	query := `
SELECT

	id,

	entity_type,

	entity_id,

	direction,

	status,

	duration_seconds,

	summary,

	follow_up_at,

	created_by,

	updated_by,

	created_at,

	updated_at

FROM call_logs

WHERE created_by=$1

AND follow_up_at IS NOT NULL

AND follow_up_at BETWEEN $2 AND $3

ORDER BY follow_up_at ASC
`

	rows, err := r.db.Query(
		ctx,
		query,
		ownerID,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var calls []entity.CallLog

	for rows.Next() {

		var call entity.CallLog

		if err := rows.Scan(
			&call.ID,
			&call.EntityType,
			&call.EntityID,
			&call.Direction,
			&call.Status,
			&call.DurationSeconds,
			&call.Summary,
			&call.FollowUpAt,
			&call.CreatedBy,
			&call.UpdatedBy,
			&call.CreatedAt,
			&call.UpdatedAt,
		); err != nil {
			return nil, err
		}

		calls = append(calls, call)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return calls, nil
}
