// Package seeders contains concrete Seeder implementations. Each implements the
// seed.Seeder interface (structurally) and must be idempotent.
package seeders

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// AdminSeeder creates a default admin user with a real bcrypt password hash.
// It replaces the previous dev_seed.sql, which stored a fake, non-bcrypt hash
// (login would have failed) and was not idempotent.
type AdminSeeder struct{}

func (AdminSeeder) Name() string { return "admin_user" }

func (AdminSeeder) Run(ctx context.Context, db *pgxpool.Pool) error {
	email := envOr("SEED_ADMIN_EMAIL", "admin@crmlite.com")
	name := envOr("SEED_ADMIN_NAME", "Admin User")
	password := envOr("SEED_ADMIN_PASSWORD", "Admin@12345")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
	`, name, email, string(hash))

	return err
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
