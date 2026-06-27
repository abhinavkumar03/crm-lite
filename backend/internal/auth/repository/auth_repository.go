package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/abhinavkumar03/crm-lite/backend/internal/auth"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(
	ctx context.Context,
	user *auth.User,
) error {

	query := `
	INSERT INTO users
	(
		id,
		name,
		email,
		password_hash
	)
	VALUES
	(
		$1,
		$2,
		$3,
		$4
	)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *AuthRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*auth.User, error) {

	query := `
	SELECT
		id,
		name,
		email,
		password_hash,
		created_at,
		updated_at
	FROM users
	WHERE email = $1
	`

	var user auth.User

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) FindByID(
	ctx context.Context,
	id string,
) (*auth.User, error) {

	query := `
	SELECT
		id,
		name,
		email,
		password_hash,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`

	var user auth.User

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
