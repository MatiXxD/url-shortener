package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connTimeout = 10 // connTimeout connection timeout in seconds
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

type DB struct {
	Pool *pgxpool.Pool
}

func (cfg *Config) ConnStr() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.SSLMode,
	)
}

func New() (*pgxpool.Pool, error) {
	cfg := readEnvs()

	pgCfg, err := pgxpool.ParseConfig(cfg.ConnStr())
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

	return pool, nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func Wait(maxWait time.Duration, pingInterval time.Duration) error {
	cfg := readEnvs()
	connStr := cfg.ConnStr()

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

func readEnvs() *Config {
	return &Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Database: os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}
}
