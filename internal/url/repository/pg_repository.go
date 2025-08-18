package repository

import (
	"context"
	"fmt"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"github.com/MatiXxD/url-shortener/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

type PostgresRepository struct {
	db     *postgres.DB
	logger *logger.Logger
}

func NewPostgresRepository(db *postgres.DB, logger *logger.Logger) *PostgresRepository {
	return &PostgresRepository{
		db:     db,
		logger: logger,
	}
}

func (pr *PostgresRepository) AddURL(ctx context.Context, url *models.URL) (string, error) {
	query := `
		INSERT INTO url (correlation_id, original, short)
		VALUES ($1, $2, $3)
		ON CONFLICT(original) DO UPDATE SET
			original = EXCLUDED.original
		RETURNING short
	`

	row := pr.db.Pool.QueryRow(ctx, query, url.CorrelationID, url.BaseURL, url.ShortURL)

	var shortURL string

	err := row.Scan(&shortURL)
	if err != nil {
		return "", fmt.Errorf("postgres add url failed with: %w", err)
	}

	return shortURL, nil
}

func (pr *PostgresRepository) BatchAddURL(ctx context.Context, urls []*models.URL) ([]*models.URL, error) {
	tx, err := pr.db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to save urls: %v", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO url (correlation_id, original, short)
		VALUES ($1, $2, $3)
		ON CONFLICT(original) DO UPDATE SET
			original = EXCLUDED.original
		RETURNING correlation_id, original, short
	`

	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue(query, url.CorrelationID, url.BaseURL, url.ShortURL)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	res := make([]*models.URL, 0, len(urls))
	for _, u := range urls {
		var url models.URL

		err := br.QueryRow().Scan(&url.CorrelationID, &url.BaseURL, &url.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to save url=%s: %w", u.BaseURL, err)
		}

		res = append(res, &url)
	}

	if err := br.Close(); err != nil {
		return nil, fmt.Errorf("failed to save urls: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to save urls: %w", err)
	}

	return res, nil
}

func (pr *PostgresRepository) GetURL(ctx context.Context, shortURL string) (*models.URL, error) {
	query := `
		SELECT correlation_id, original, short FROM url
		WHERE short = $1
	`

	row := pr.db.Pool.QueryRow(ctx, query, shortURL)

	var url models.URL

	err := row.Scan(&url.CorrelationID, &url.BaseURL, &url.ShortURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	return &url, nil
}
