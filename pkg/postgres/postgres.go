package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connTimeout = 10 // connTimeout connection timeout in seconds
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(connStr string) (*DB, error) {
	pgCfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres config: %w", err)
	}

	// Disable statement cache for PgBouncer compatibility
	pgCfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	ctx, cancel := context.WithTimeout(context.Background(), connTimeout*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, pgCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func Wait(connStr string, maxWait time.Duration, pingInterval time.Duration) error {
	deadline := time.Now().Add(maxWait)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), pingInterval)
		pool, err := pgxpool.New(ctx, connStr)

		if err == nil {
			err = pool.Ping(ctx)
			pool.Close()

			if err == nil {
				cancel()
				return nil
			}
		}

		cancel()

		if time.Now().After(deadline) {
			return fmt.Errorf("deadline exceeded: %w", err)
		}

		time.Sleep(pingInterval)
	}
}
