package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps the pgx connection pool.
type DB struct {
	pool *pgxpool.Pool
}

// Connect creates a new database connection pool using the DATABASE_URL environment variable.
func Connect() (*DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to reach database: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Pool returns the underlying pgx pool for use by other packages.
func (d *DB) Pool() *pgxpool.Pool {
	return d.pool
}

// Close shuts down the connection pool.
func (d *DB) Close() {
	d.pool.Close()
}
