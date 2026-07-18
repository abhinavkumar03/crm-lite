// Package seed provides a small, ordered, history-tracked seeder framework.
//
// Design:
//   - Each seed is a Seeder with a stable Name() used as its identity.
//   - The Runner records applied seeders in a schema_seed_history table so
//     seeders run at most once (unless run with fresh=true).
//   - Seeders should additionally be idempotent (e.g. INSERT ... ON CONFLICT)
//     so a fresh run is always safe.
//
// This keeps seeding versioned and repeatable across dev/CI without duplicating
// data, while remaining simple enough to read end-to-end.
package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Seeder populates the database with a specific set of data.
type Seeder interface {
	// Name is the stable identity of the seeder (used for history tracking).
	Name() string
	// Run inserts the seeder's data. Implementations must be idempotent.
	Run(ctx context.Context, db *pgxpool.Pool) error
}
