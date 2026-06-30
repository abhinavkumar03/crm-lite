package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/note/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/note/entity"
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
	note *entity.Note,
) error {

	query := `
	INSERT INTO notes (
		id,
		entity_type,
		entity_id,
		note,
		created_by,
		updated_by,
		created_at,
		updated_at
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8
	)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		note.ID,
		note.EntityType,
		note.EntityID,
		note.Note,
		note.CreatedBy,
		note.UpdatedBy,
		note.CreatedAt,
		note.UpdatedAt,
	)

	return err
}

func (r *Repository) List(
	ctx context.Context,
	entityType string,
	entityID string,
) ([]dto.NoteResponse, error) {

	query := `
	SELECT
		n.id,
		n.entity_type,
		n.entity_id,
		n.note,
		n.created_by,
		u.name,
		n.updated_by,
		n.created_at,
		n.updated_at
	FROM notes n
	INNER JOIN users u
		ON u.id = n.created_by
	WHERE n.entity_type = $1
	AND n.entity_id = $2
	ORDER BY n.created_at DESC;
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

	notes := make([]dto.NoteResponse, 0)

	for rows.Next() {

		var note dto.NoteResponse

		err := rows.Scan(
			&note.ID,
			&note.EntityType,
			&note.EntityID,
			&note.Note,
			&note.CreatedBy,
			&note.User.Name,
			&note.UpdatedBy,
			&note.CreatedAt,
			&note.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		note.User.ID = note.CreatedBy

		notes = append(
			notes,
			note,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*entity.Note, error) {

	query := `
	SELECT
		id,
		entity_type,
		entity_id,
		note,
		created_by,
		updated_by,
		created_at,
		updated_at
	FROM notes
	WHERE id=$1
	`

	var note entity.Note

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&note.ID,
		&note.EntityType,
		&note.EntityID,
		&note.Note,
		&note.CreatedBy,
		&note.UpdatedBy,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *Repository) Update(
	ctx context.Context,
	note *entity.Note,
) error {

	query := `
	UPDATE notes
	SET
		note=$1,
		updated_by=$2,
		updated_at=$3
	WHERE id=$4
	`

	_, err := r.db.Exec(
		ctx,
		query,
		note.Note,
		note.UpdatedBy,
		note.UpdatedAt,
		note.ID,
	)

	return err
}

func (r *Repository) Delete(
	ctx context.Context,
	id string,
) error {

	query := `
	DELETE
	FROM notes
	WHERE id=$1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}
