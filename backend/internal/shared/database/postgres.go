package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(connectionString string) (*pgxpool.Pool, error) {
	return pgxpool.New(
		context.Background(),
		connectionString,
	)
}

func BuildDSN(
	host,
	port,
	user,
	password,
	dbName,
	sslMode string,
) string {

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		dbName,
		sslMode,
	)
}
