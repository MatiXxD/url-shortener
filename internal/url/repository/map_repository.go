package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/pkg/logger"
)

type MapRepository struct {
	db     map[string]*models.URL
	pk     int
	logger *logger.Logger
	mu     sync.RWMutex
}

func NewMapRepository(d map[string]*models.URL, l *logger.Logger) *MapRepository {
	return &MapRepository{
		db:     d,
		pk:     1,
		logger: l,
		mu:     sync.RWMutex{},
	}
}

func (mr *MapRepository) AddURL(ctx context.Context, shortenURL *models.URL) (string, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if got, ok := mr.db[shortenURL.BaseURL]; ok {
		return got.ShortURL, nil
	}

	mr.db[shortenURL.BaseURL] = &models.URL{
		ID:            mr.pk,
		CorrelationID: shortenURL.CorrelationID,
		BaseURL:       shortenURL.BaseURL,
		ShortURL:      shortenURL.ShortURL,
		CreateAt:      time.Now(),
	}
	mr.pk++

	return shortenURL.ShortURL, nil
}

func (mr *MapRepository) BatchAddURL(ctx context.Context, urls []*models.URL) ([]*models.URL, error) {
	res := make([]*models.URL, 0, len(urls))

	for _, u := range urls {
		shortUrl, err := mr.AddURL(ctx, u)
		if err != nil {
			return res, fmt.Errorf("failed to add url=%s: %w", u.BaseURL, err)
		}

		res = append(res, &models.URL{
			CorrelationID: u.CorrelationID,
			BaseURL:       u.BaseURL,
			ShortURL:      shortUrl,
		})
	}

	return res, nil
}

func (mr *MapRepository) GetURL(ctx context.Context, shortURL string) (*models.URL, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	for _, v := range mr.db {
		if v.ShortURL == shortURL {
			return &models.URL{
				CorrelationID: v.CorrelationID,
				BaseURL:       v.BaseURL,
				ShortURL:      v.ShortURL,
			}, nil
		}
	}

	return nil, fmt.Errorf("url was not found")
}
