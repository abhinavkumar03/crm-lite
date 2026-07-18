package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/validationengine/entity"
)

const ruleColumns = `
	id, organization_id, module_id, field_id, rule_type,
	params, error_message, is_active, sort_order, created_at, updated_at
`

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func scanRule(row pgx.Row, r *entity.ValidationRule) error {
	return row.Scan(
		&r.ID, &r.OrganizationID, &r.ModuleID, &r.FieldID, &r.RuleType,
		&r.Params, &r.ErrorMessage, &r.IsActive, &r.SortOrder, &r.CreatedAt, &r.UpdatedAt,
	)
}

func (r *Repository) ModuleExists(ctx context.Context, orgID, moduleID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM modules WHERE id = $1 AND organization_id = $2)
	`, moduleID, orgID).Scan(&exists)
	return exists, err
}

// FieldExists reports whether a field belongs to the given module (used to
// validate that a rule's field_id targets a field on the same module).
func (r *Repository) FieldExists(ctx context.Context, orgID, moduleID, fieldID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM fields
			WHERE id = $1 AND module_id = $2 AND organization_id = $3
		)
	`, fieldID, moduleID, orgID).Scan(&exists)
	return exists, err
}

func (r *Repository) Create(ctx context.Context, rule *entity.ValidationRule) error {
	return r.db.QueryRow(ctx, `
		INSERT INTO validation_rules (
			organization_id, module_id, field_id, rule_type,
			params, error_message, is_active, sort_order
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at, updated_at
	`,
		rule.OrganizationID, rule.ModuleID, rule.FieldID, rule.RuleType,
		rule.Params, rule.ErrorMessage, rule.IsActive, rule.SortOrder,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

func (r *Repository) List(ctx context.Context, orgID, moduleID string) ([]entity.ValidationRule, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+ruleColumns+`
		FROM validation_rules
		WHERE organization_id = $1 AND module_id = $2
		ORDER BY sort_order ASC, created_at ASC
	`, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collect(rows)
}

// ActiveByModule returns only active rules, used by the engine at evaluation time.
func (r *Repository) ActiveByModule(ctx context.Context, orgID, moduleID string) ([]entity.ValidationRule, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+ruleColumns+`
		FROM validation_rules
		WHERE organization_id = $1 AND module_id = $2 AND is_active = TRUE
		ORDER BY sort_order ASC, created_at ASC
	`, orgID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collect(rows)
}

func (r *Repository) GetByID(ctx context.Context, orgID, moduleID, id string) (*entity.ValidationRule, error) {
	var rule entity.ValidationRule
	err := scanRule(r.db.QueryRow(ctx, `
		SELECT `+ruleColumns+`
		FROM validation_rules
		WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID), &rule)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *Repository) Update(ctx context.Context, rule *entity.ValidationRule) error {
	return r.db.QueryRow(ctx, `
		UPDATE validation_rules SET
			params = $1,
			error_message = $2,
			is_active = $3,
			sort_order = $4,
			updated_at = NOW()
		WHERE id = $5 AND module_id = $6 AND organization_id = $7
		RETURNING updated_at
	`,
		rule.Params, rule.ErrorMessage, rule.IsActive, rule.SortOrder,
		rule.ID, rule.ModuleID, rule.OrganizationID,
	).Scan(&rule.UpdatedAt)
}

func (r *Repository) Delete(ctx context.Context, orgID, moduleID, id string) (bool, error) {
	tag, err := r.db.Exec(ctx, `
		DELETE FROM validation_rules WHERE id = $1 AND module_id = $2 AND organization_id = $3
	`, id, moduleID, orgID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func collect(rows pgx.Rows) ([]entity.ValidationRule, error) {
	rules := make([]entity.ValidationRule, 0)
	for rows.Next() {
		var rule entity.ValidationRule
		if err := scanRule(rows, &rule); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, rows.Err()
}
