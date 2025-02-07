package repository

import (
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/pkg/tokengen"
)

const token_size = 10

type MapRepository struct {
	db map[string]*models.URL
}

func NewMapRepository(d map[string]*models.URL) url.Repository {
	return &MapRepository{
		db: d,
	}
}

func (mr *MapRepository) ReduceURL(url string) (string, error) {
	if url, ok := mr.db[url]; ok {
		return url.ShortURL, nil
	}

	shortURL := tokengen.GenerateToken(token_size)
	mr.db[url] = models.NewURL(url, shortURL)

	return shortURL, nil
}

func (mr *MapRepository) GetURL(shortURL string) (string, bool) {
	for k, v := range mr.db {
		if v.ShortURL == shortURL {
			return k, true
		}
	}
	return "", false
}
