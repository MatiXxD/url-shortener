package repository

import (
	"sync"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/pkg/logger"
)

type MapRepository struct {
	db     map[string]*models.URL
	logger *logger.Logger
	mu     sync.RWMutex
}

func NewMapRepository(d map[string]*models.URL, l *logger.Logger) url.Repository {
	return &MapRepository{
		db:     d,
		logger: l,
		mu:     sync.RWMutex{},
	}
}

func (mr *MapRepository) AddURL(url, shortURL string) (string, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	if got, ok := mr.db[url]; ok {
		return got.ShortURL, nil
	}
	mr.db[url] = models.NewURL(url, shortURL)
	return shortURL, nil
}

func (mr *MapRepository) GetURL(shortURL string) (string, bool) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	for k, v := range mr.db {
		if v.ShortURL == shortURL {
			return k, true
		}
	}
	return "", false
}
