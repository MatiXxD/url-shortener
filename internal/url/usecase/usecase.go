package usecase

import (
	"fmt"

	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"github.com/MatiXxD/url-shortener/pkg/tokengen"
)

const tokenSize = 10

type UrlUsecase struct {
	repo   url.Repository
	logger *logger.Logger
}

func NewUrlUsecase(r url.Repository, l *logger.Logger) url.Usecase {
	return &UrlUsecase{
		repo:   r,
		logger: l,
	}
}

func (uu *UrlUsecase) ReduceURL(url string) (string, error) {
	genURL := tokengen.GenerateToken(tokenSize)
	shortURL, err := uu.repo.AddURL(url, genURL)
	if err != nil {
		uu.logger.Error("can't add short url to database")
		return "", fmt.Errorf("can't add short url to database: %v", err)
	}
	return shortURL, nil
}

func (uu *UrlUsecase) GetURL(shortURL string) (string, bool) {
	url, ok := uu.repo.GetURL(shortURL)
	if !ok {
		return "", false
	}
	return url, true
}
