package seed

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// Runner executes registered seeders in order, skipping ones already recorded
// in schema_seed_history.
type Runner struct {
	db      *pgxpool.Pool
	logger  *zap.Logger
	seeders []Seeder
}

func NewRunner(db *pgxpool.Pool, logger *zap.Logger) *Runner {
	return &Runner{db: db, logger: logger}
}

// Register appends seeders to the run list. Order of registration is the order
// of execution, so dependencies (e.g. users before leads) must be registered
// first.
func (r *Runner) Register(seeders ...Seeder) {
	r.seeders = append(r.seeders, seeders...)
}

// Run applies pending seeders. When fresh is true, the history is cleared first
// so every seeder runs again (safe because seeders are idempotent).
func (r *Runner) Run(ctx context.Context, fresh bool) error {
	if err := r.ensureHistory(ctx); err != nil {
		return err
	}

	if fresh {
		if _, err := r.db.Exec(ctx, `DELETE FROM schema_seed_history`); err != nil {
			return fmt.Errorf("seed: clear history: %w", err)
		}
		r.logger.Info("seed: history cleared (fresh run)")
	}

	for _, s := range r.seeders {
		done, err := r.isApplied(ctx, s.Name())
		if err != nil {
			return err
		}
		if done {
			r.logger.Info("seed: skip (already applied)", zap.String("seeder", s.Name()))
			continue
		}

		r.logger.Info("seed: running", zap.String("seeder", s.Name()))
		if err := s.Run(ctx, r.db); err != nil {
			return fmt.Errorf("seed %q: %w", s.Name(), err)
		}

		if err := r.markApplied(ctx, s.Name()); err != nil {
			return err
		}
		r.logger.Info("seed: applied", zap.String("seeder", s.Name()))
	}

	return nil
}

func (r *Runner) ensureHistory(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_seed_history (
			name        TEXT PRIMARY KEY,
			applied_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("seed: ensure history table: %w", err)
	}
	return nil
}

func (r *Runner) isApplied(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM schema_seed_history WHERE name = $1)
	`, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("seed: check history: %w", err)
	}
	return exists, nil
}

func (r *Runner) markApplied(ctx context.Context, name string) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO schema_seed_history (name) VALUES ($1)
		ON CONFLICT (name) DO NOTHING
	`, name)
	if err != nil {
		return fmt.Errorf("seed: record history: %w", err)
	}
	return nil
}
