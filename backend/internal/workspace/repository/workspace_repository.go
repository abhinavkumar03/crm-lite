package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const EntityRecord = "RECORD"

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) RecordExists(ctx context.Context, orgID, moduleID, recordID string) (bool, error) {
	var ok bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM records
			WHERE id = $1 AND organization_id = $2 AND module_id = $3
		)
	`, recordID, orgID, moduleID).Scan(&ok)
	return ok, err
}

// --- Layouts ----------------------------------------------------------------

type Layout struct {
	ID       string
	Name     string
	Type     string
	IsDefault bool
	Config   json.RawMessage
}

func (r *Repository) GetDefaultDetailLayout(ctx context.Context, orgID, moduleID string) (*Layout, error) {
	var l Layout
	err := r.db.QueryRow(ctx, `
		SELECT id::text, name, layout_type, is_default, config
		FROM layouts
		WHERE organization_id = $1 AND module_id = $2
		  AND layout_type = 'detail' AND is_default = TRUE
		LIMIT 1
	`, orgID, moduleID).Scan(&l.ID, &l.Name, &l.Type, &l.IsDefault, &l.Config)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repository) UpsertDefaultDetailLayout(ctx context.Context, orgID, moduleID string, config json.RawMessage) (*Layout, error) {
	existing, err := r.GetDefaultDetailLayout(ctx, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		_, err = r.db.Exec(ctx, `
			UPDATE layouts SET config = $2, updated_at = NOW()
			WHERE id = $1
		`, existing.ID, config)
		if err != nil {
			return nil, err
		}
		existing.Config = config
		return existing, nil
	}
	var l Layout
	err = r.db.QueryRow(ctx, `
		INSERT INTO layouts (organization_id, module_id, name, layout_type, is_default, config)
		VALUES ($1, $2, 'Default Detail', 'detail', TRUE, $3)
		RETURNING id::text, name, layout_type, is_default, config
	`, orgID, moduleID, config).Scan(&l.ID, &l.Name, &l.Type, &l.IsDefault, &l.Config)
	return &l, err
}

func (r *Repository) ListVisibleFieldAPINames(ctx context.Context, orgID, moduleID string) ([]string, error) {
	rows, err := r.db.Query(ctx, `
		SELECT api_name FROM fields
		WHERE organization_id = $1 AND module_id = $2 AND is_visible = TRUE
		ORDER BY sort_order, api_name
	`, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		out = append(out, name)
	}
	return out, rows.Err()
}

// --- Notes ------------------------------------------------------------------

type Note struct {
	ID        string
	Title     *string
	Body      string
	CreatedBy string
	AuthorName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Repository) CreateNote(ctx context.Context, orgID, moduleID, recordID, userID, body string, title *string) (*Note, error) {
	id := uuid.NewString()
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx, `
		INSERT INTO notes (
			id, entity_type, entity_id, note, title,
			organization_id, module_id, created_by, updated_by, created_at, updated_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$7,$8,$8)
	`, id, recordID, body, title, orgID, moduleID, userID, now)
	if err != nil {
		return nil, err
	}
	return r.GetNote(ctx, orgID, id)
}

func (r *Repository) GetNote(ctx context.Context, orgID, noteID string) (*Note, error) {
	var n Note
	err := r.db.QueryRow(ctx, `
		SELECT n.id::text, n.title, n.note, n.created_by::text, u.name, n.created_at, n.updated_at
		FROM notes n
		JOIN users u ON u.id = n.created_by
		WHERE n.id = $1 AND n.organization_id = $2 AND n.entity_type = 'RECORD'
	`, noteID, orgID).Scan(&n.ID, &n.Title, &n.Body, &n.CreatedBy, &n.AuthorName, &n.CreatedAt, &n.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &n, err
}

func (r *Repository) ListNotes(ctx context.Context, orgID, moduleID, recordID string) ([]Note, error) {
	rows, err := r.db.Query(ctx, `
		SELECT n.id::text, n.title, n.note, n.created_by::text, u.name, n.created_at, n.updated_at
		FROM notes n
		JOIN users u ON u.id = n.created_by
		WHERE n.organization_id = $1 AND n.module_id = $2
		  AND n.entity_type = 'RECORD' AND n.entity_id = $3
		ORDER BY n.created_at DESC
	`, orgID, moduleID, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedBy, &n.AuthorName, &n.CreatedAt, &n.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *Repository) DeleteNote(ctx context.Context, orgID, noteID string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM notes WHERE id = $1 AND organization_id = $2 AND entity_type = 'RECORD'
	`, noteID, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

// --- Attachments ------------------------------------------------------------

type Attachment struct {
	ID           string
	FileName     string
	FileURL      string
	PublicID     string
	ResourceType *string
	FileSize     *int64
	UploadedBy   string
	UploaderName string
	CreatedAt    time.Time
}

func (r *Repository) CreateAttachment(
	ctx context.Context,
	orgID, moduleID, recordID, userID string,
	fileName, fileURL, publicID string,
	resourceType *string, fileSize *int64,
) (*Attachment, error) {
	id := uuid.NewString()
	now := time.Now().UTC()
	_, err := r.db.Exec(ctx, `
		INSERT INTO attachments (
			id, entity_type, entity_id, file_name, file_url, public_id,
			resource_type, file_size, uploaded_by, organization_id, module_id, created_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`, id, recordID, fileName, fileURL, publicID, resourceType, fileSize, userID, orgID, moduleID, now)
	if err != nil {
		return nil, err
	}
	return r.GetAttachment(ctx, orgID, id)
}

func (r *Repository) GetAttachment(ctx context.Context, orgID, id string) (*Attachment, error) {
	var a Attachment
	err := r.db.QueryRow(ctx, `
		SELECT a.id::text, a.file_name, a.file_url, a.public_id, a.resource_type, a.file_size,
		       a.uploaded_by::text, u.name, a.created_at
		FROM attachments a
		JOIN users u ON u.id = a.uploaded_by
		WHERE a.id = $1 AND a.organization_id = $2 AND a.entity_type = 'RECORD'
	`, id, orgID).Scan(
		&a.ID, &a.FileName, &a.FileURL, &a.PublicID, &a.ResourceType, &a.FileSize,
		&a.UploadedBy, &a.UploaderName, &a.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &a, err
}

func (r *Repository) ListAttachments(ctx context.Context, orgID, moduleID, recordID string) ([]Attachment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id::text, a.file_name, a.file_url, a.public_id, a.resource_type, a.file_size,
		       a.uploaded_by::text, u.name, a.created_at
		FROM attachments a
		JOIN users u ON u.id = a.uploaded_by
		WHERE a.organization_id = $1 AND a.module_id = $2
		  AND a.entity_type = 'RECORD' AND a.entity_id = $3
		ORDER BY a.created_at DESC
	`, orgID, moduleID, recordID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Attachment, 0)
	for rows.Next() {
		var a Attachment
		if err := rows.Scan(
			&a.ID, &a.FileName, &a.FileURL, &a.PublicID, &a.ResourceType, &a.FileSize,
			&a.UploadedBy, &a.UploaderName, &a.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *Repository) DeleteAttachment(ctx context.Context, orgID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM attachments WHERE id = $1 AND organization_id = $2 AND entity_type = 'RECORD'
	`, id, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

// --- Activities -------------------------------------------------------------

type Activity struct {
	ID          string
	Action      string
	Description string
	PerformedBy string
	ActorName   string
	Metadata    json.RawMessage
	CreatedAt   time.Time
}

func (r *Repository) CreateActivity(
	ctx context.Context,
	orgID, moduleID, recordID, userID, action, description string,
	metadata json.RawMessage,
) error {
	if metadata == nil {
		metadata = json.RawMessage(`{}`)
	}
	_, err := r.db.Exec(ctx, `
		INSERT INTO activities (
			id, entity_type, entity_id, action, description, performed_by,
			metadata, organization_id, module_id, created_at
		) VALUES ($1,'RECORD',$2,$3,$4,$5,$6,$7,$8,NOW())
	`, uuid.NewString(), recordID, action, description, userID, metadata, orgID, moduleID)
	return err
}

func (r *Repository) ListActivities(ctx context.Context, orgID, moduleID, recordID string, limit int) ([]Activity, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.db.Query(ctx, `
		SELECT a.id::text, a.action, a.description, a.performed_by::text, u.name,
		       COALESCE(a.metadata, '{}'::jsonb), a.created_at
		FROM activities a
		JOIN users u ON u.id = a.performed_by
		WHERE a.organization_id = $1 AND a.module_id = $2
		  AND a.entity_type = 'RECORD' AND a.entity_id = $3
		ORDER BY a.created_at DESC
		LIMIT $4
	`, orgID, moduleID, recordID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Activity, 0)
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.Action, &a.Description, &a.PerformedBy, &a.ActorName, &a.Metadata, &a.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// --- Related lists ----------------------------------------------------------

type RelatedDescriptor struct {
	ChildModuleID   string
	ChildModuleName string
	ChildAPIName    string
	LookupFieldAPI  string
	LookupFieldLabel string
}

func (r *Repository) ListRelatedDescriptors(ctx context.Context, orgID, parentModuleID string) ([]RelatedDescriptor, error) {
	rows, err := r.db.Query(ctx, `
		SELECT m.id::text, m.plural_label, m.api_name, f.api_name, f.label
		FROM fields f
		JOIN modules m ON m.id = f.module_id
		WHERE f.organization_id = $1
		  AND f.field_type = 'lookup'
		  AND f.lookup_module_id = $2
		  AND m.is_enabled = TRUE
		ORDER BY m.sort_order, f.sort_order
	`, orgID, parentModuleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]RelatedDescriptor, 0)
	for rows.Next() {
		var d RelatedDescriptor
		if err := rows.Scan(&d.ChildModuleID, &d.ChildModuleName, &d.ChildAPIName, &d.LookupFieldAPI, &d.LookupFieldLabel); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *Repository) LookupFieldAPI(ctx context.Context, orgID, childModuleID, parentModuleID string) (string, error) {
	var api string
	err := r.db.QueryRow(ctx, `
		SELECT api_name FROM fields
		WHERE organization_id = $1 AND module_id = $2
		  AND field_type = 'lookup' AND lookup_module_id = $3
		ORDER BY sort_order
		LIMIT 1
	`, orgID, childModuleID, parentModuleID).Scan(&api)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("no lookup relationship")
	}
	return api, err
}
