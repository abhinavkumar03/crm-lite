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

func (r *AuthRepository) CreateUser(
	ctx context.Context,
	user *auth.User,
) error {

	query := `
	INSERT INTO users
	(
		name,
		email,
		password_hash
	)
	VALUES
	(
		$1,
		$2,
		$3
	)
	RETURNING
		id,
		created_at,
		updated_at;
	`

	return r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
}

func (r *AuthRepository) GetUserByEmail(
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
	WHERE email = $1;
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

func (r *AuthRepository) GetUserByID(
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
	WHERE id = $1;
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

func (r *AuthRepository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE email = $1
	)
	`

	var exists bool

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
