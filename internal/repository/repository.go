package repository

import (
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/pkg/tokengen"
)

const token_size = 10

var mockDB map[string]*models.URL = make(map[string]*models.URL)

func ReduceURL(url string) (string, error) {
	if url, ok := mockDB[url]; ok {
		return url.ShortURL, nil
	}

	shortURL := tokengen.GenerateToken(token_size)
	mockDB[url] = models.NewURL(url, shortURL)

	return shortURL, nil
}

func GetURL(shortURL string) (string, bool) {
	for k, v := range mockDB {
		if v.ShortURL == shortURL {
			return k, true
		}
	}
	return "", false
}
