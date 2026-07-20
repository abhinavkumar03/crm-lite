package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/notification/entity"
)

type ProviderRepository struct {
	db *pgxpool.Pool
}

func NewProviderRepository(db *pgxpool.Pool) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) Create(ctx context.Context, p *entity.Provider) error {
	if p.Config == nil {
		p.Config = []byte("{}")
	}
	return r.db.QueryRow(ctx, `
		INSERT INTO communication_providers (
			organization_id, channel, provider_type, name, config, secrets_encrypted,
			is_default, is_active, created_by
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at
	`,
		p.OrganizationID, p.Channel, p.ProviderType, p.Name, p.Config, p.SecretsEncrypted,
		p.IsDefault, p.IsActive, p.CreatedBy,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *ProviderRepository) Update(ctx context.Context, p *entity.Provider) error {
	_, err := r.db.Exec(ctx, `
		UPDATE communication_providers SET
			name = $3, provider_type = $4, config = $5,
			secrets_encrypted = COALESCE($6, secrets_encrypted),
			is_default = $7, is_active = $8, updated_at = NOW()
		WHERE id = $1 AND organization_id = $2
	`, p.ID, p.OrganizationID, p.Name, p.ProviderType, p.Config, p.SecretsEncrypted, p.IsDefault, p.IsActive)
	return err
}

func (r *ProviderRepository) Delete(ctx context.Context, orgID, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM communication_providers WHERE id = $1 AND organization_id = $2`, id, orgID)
	return err
}

func (r *ProviderRepository) Get(ctx context.Context, orgID, id string) (*entity.Provider, error) {
	var p entity.Provider
	err := r.db.QueryRow(ctx, `
		SELECT id, organization_id, channel, provider_type, name, config, secrets_encrypted,
		       is_default, is_active, last_health_at, last_error, created_by, created_at, updated_at
		FROM communication_providers WHERE id = $1 AND organization_id = $2
	`, id, orgID).Scan(
		&p.ID, &p.OrganizationID, &p.Channel, &p.ProviderType, &p.Name, &p.Config, &p.SecretsEncrypted,
		&p.IsDefault, &p.IsActive, &p.LastHealthAt, &p.LastError, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.SecretsConfigured = len(p.SecretsEncrypted) > 0
	return &p, nil
}

func (r *ProviderRepository) List(ctx context.Context, orgID, channel string) ([]entity.Provider, error) {
	q := `
		SELECT id, organization_id, channel, provider_type, name, config, secrets_encrypted,
		       is_default, is_active, last_health_at, last_error, created_by, created_at, updated_at
		FROM communication_providers WHERE organization_id = $1`
	args := []any{orgID}
	if channel != "" {
		q += ` AND channel = $2`
		args = append(args, channel)
	}
	q += ` ORDER BY is_default DESC, name ASC`

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Provider, 0)
	for rows.Next() {
		var p entity.Provider
		if err := rows.Scan(
			&p.ID, &p.OrganizationID, &p.Channel, &p.ProviderType, &p.Name, &p.Config, &p.SecretsEncrypted,
			&p.IsDefault, &p.IsActive, &p.LastHealthAt, &p.LastError, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.SecretsConfigured = len(p.SecretsEncrypted) > 0
		items = append(items, p)
	}
	return items, rows.Err()
}

func (r *ProviderRepository) GetDefault(ctx context.Context, orgID, channel string) (*entity.Provider, error) {
	var p entity.Provider
	err := r.db.QueryRow(ctx, `
		SELECT id, organization_id, channel, provider_type, name, config, secrets_encrypted,
		       is_default, is_active, last_health_at, last_error, created_by, created_at, updated_at
		FROM communication_providers
		WHERE organization_id = $1 AND channel = $2 AND is_active = TRUE AND is_default = TRUE
		LIMIT 1
	`, orgID, channel).Scan(
		&p.ID, &p.OrganizationID, &p.Channel, &p.ProviderType, &p.Name, &p.Config, &p.SecretsEncrypted,
		&p.IsDefault, &p.IsActive, &p.LastHealthAt, &p.LastError, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		// Fall back to any active provider for the channel.
		err = r.db.QueryRow(ctx, `
			SELECT id, organization_id, channel, provider_type, name, config, secrets_encrypted,
			       is_default, is_active, last_health_at, last_error, created_by, created_at, updated_at
			FROM communication_providers
			WHERE organization_id = $1 AND channel = $2 AND is_active = TRUE
			ORDER BY is_default DESC, updated_at DESC
			LIMIT 1
		`, orgID, channel).Scan(
			&p.ID, &p.OrganizationID, &p.Channel, &p.ProviderType, &p.Name, &p.Config, &p.SecretsEncrypted,
			&p.IsDefault, &p.IsActive, &p.LastHealthAt, &p.LastError, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}
	if err != nil {
		return nil, err
	}
	p.SecretsConfigured = len(p.SecretsEncrypted) > 0
	return &p, nil
}

func (r *ProviderRepository) ClearDefault(ctx context.Context, orgID, channel string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE communication_providers SET is_default = FALSE, updated_at = NOW()
		WHERE organization_id = $1 AND channel = $2 AND is_default = TRUE
	`, orgID, channel)
	return err
}

func (r *ProviderRepository) MarkHealth(ctx context.Context, id string, errMsg *string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE communication_providers
		SET last_health_at = NOW(), last_error = $2, updated_at = NOW()
		WHERE id = $1
	`, id, errMsg)
	return err
}

func (r *ProviderRepository) CreateSender(ctx context.Context, s *entity.SenderIdentity) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO communication_sender_identities (
			organization_id, provider_id, channel, display_name, from_address, reply_to, is_default
		) VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at
	`, s.OrganizationID, s.ProviderID, s.Channel, s.DisplayName, s.FromAddress, s.ReplyTo, s.IsDefault).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *ProviderRepository) ListSenders(ctx context.Context, orgID, channel string) ([]entity.SenderIdentity, error) {
	q := `
		SELECT id, organization_id, provider_id, channel, display_name, from_address, reply_to,
		       is_default, created_at, updated_at
		FROM communication_sender_identities WHERE organization_id = $1`
	args := []any{orgID}
	if channel != "" {
		q += ` AND channel = $2`
		args = append(args, channel)
	}
	q += ` ORDER BY is_default DESC, from_address ASC`
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]entity.SenderIdentity, 0)
	for rows.Next() {
		var s entity.SenderIdentity
		if err := rows.Scan(
			&s.ID, &s.OrganizationID, &s.ProviderID, &s.Channel, &s.DisplayName,
			&s.FromAddress, &s.ReplyTo, &s.IsDefault, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

func (r *ProviderRepository) DeleteSender(ctx context.Context, orgID, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM communication_sender_identities WHERE id = $1 AND organization_id = $2`, id, orgID)
	return err
}
