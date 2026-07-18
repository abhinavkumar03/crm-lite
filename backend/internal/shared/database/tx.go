package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager provides a small Unit-of-Work abstraction over pgx so that services
// can run multi-statement operations atomically without leaking transaction
// lifecycle management (Begin/Commit/Rollback) into business logic.
//
// Repositories that need to participate in a transaction should accept a
// pgx.Tx (which shares the Query/Exec interface with the pool) instead of the
// pool directly.
type TxManager struct {
	pool *pgxpool.Pool
}

func NewTxManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{pool: pool}
}

// WithTx runs fn inside a database transaction. The transaction is committed if
// fn returns nil, rolled back if fn returns an error, and rolled back before
// re-panicking if fn panics. The named return value is used so the deferred
// commit/rollback decision can observe fn's error.
func (m *TxManager) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) (err error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("database: begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}

		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				err = fmt.Errorf("database: rollback failed: %v (original error: %w)", rbErr, err)
			}
			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			err = fmt.Errorf("database: commit transaction: %w", commitErr)
		}
	}()

	err = fn(tx)
	return err
}
