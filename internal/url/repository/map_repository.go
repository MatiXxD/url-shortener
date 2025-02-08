package repository

import (
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url"
)

type MapRepository struct {
	db map[string]*models.URL
}

func NewMapRepository(d map[string]*models.URL) url.Repository {
	return &MapRepository{
		db: d,
	}
}

func (mr *MapRepository) AddURL(url, shortURL string) (string, error) {
	if got, ok := mr.db[url]; ok {
		return got.ShortURL, nil
	}

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
