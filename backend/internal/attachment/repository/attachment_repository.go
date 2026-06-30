package repository

import (
	"context"

	"github.com/abhinavkumar03/crm-lite/backend/internal/attachment/entity"
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
	attachment *entity.Attachment,
) error {

	query := `
	INSERT INTO attachments (

		id,

		entity_type,

		entity_id,

		file_name,

		file_url,

		public_id,

		resource_type,

		file_size,

		uploaded_by,

		created_at

	)

	VALUES (

		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10

	)
	`

	_, err := r.db.Exec(
		ctx,
		query,

		attachment.ID,
		attachment.EntityType,
		attachment.EntityID,
		attachment.FileName,
		attachment.FileURL,
		attachment.PublicID,
		attachment.ResourceType,
		attachment.FileSize,
		attachment.UploadedBy,
		attachment.CreatedAt,
	)

	return err
}

func (r *Repository) List(
	ctx context.Context,
	entityType string,
	entityID string,
) ([]entity.Attachment, error) {

	query := `
	SELECT

		id,

		entity_type,

		entity_id,

		file_name,

		file_url,

		public_id,

		resource_type,

		file_size,

		uploaded_by,

		created_at

	FROM attachments

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

	attachments := make([]entity.Attachment, 0)

	for rows.Next() {

		var attachment entity.Attachment

		err := rows.Scan(
			&attachment.ID,
			&attachment.EntityType,
			&attachment.EntityID,
			&attachment.FileName,
			&attachment.FileURL,
			&attachment.PublicID,
			&attachment.ResourceType,
			&attachment.FileSize,
			&attachment.UploadedBy,
			&attachment.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		attachments = append(
			attachments,
			attachment,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attachments, nil
}

func (r *Repository) GetByID(
	ctx context.Context,
	id string,
) (*entity.Attachment, error) {

	query := `
	SELECT

		id,

		entity_type,

		entity_id,

		file_name,

		file_url,

		public_id,

		resource_type,

		file_size,

		uploaded_by,

		created_at

	FROM attachments

	WHERE id=$1
	`

	var attachment entity.Attachment

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&attachment.ID,
		&attachment.EntityType,
		&attachment.EntityID,
		&attachment.FileName,
		&attachment.FileURL,
		&attachment.PublicID,
		&attachment.ResourceType,
		&attachment.FileSize,
		&attachment.UploadedBy,
		&attachment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &attachment, nil
}
func (r *Repository) Delete(
	ctx context.Context,
	id string,
) error {

	query := `
	DELETE
	FROM attachments
	WHERE id=$1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		id,
	)

	return err
}
