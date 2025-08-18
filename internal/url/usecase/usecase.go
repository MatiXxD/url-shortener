package usecase

import (
	"context"
	"fmt"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"github.com/MatiXxD/url-shortener/pkg/tokengen"
)

const (
	tokenSize = 10
	batchSize = 5
)

type UrlUsecase struct {
	repo   url.Repository
	cfg    *config.ServiceConfig
	logger *logger.Logger
}

func NewUrlUsecase(r url.Repository, cfg *config.ServiceConfig, l *logger.Logger) *UrlUsecase {
	return &UrlUsecase{
		repo:   r,
		cfg:    cfg,
		logger: l,
	}
}

func (uu *UrlUsecase) ReduceURL(ctx context.Context, req *models.UrlDTO) (string, error) {
	genURL := tokengen.GenerateToken(tokenSize)

	shortURL, err := uu.repo.AddURL(ctx, &models.URL{
		CorrelationID: req.CorrelationID,
		BaseURL:       req.OriginURL,
		ShortURL:      genURL,
	})
	if err != nil {
		uu.logger.Error("can't add short url to database")
		return "", fmt.Errorf("can't add short url to database: %v", err)
	}

	return uu.getShortURL(shortURL), nil
}

func (uu *UrlUsecase) BatchReduceURL(ctx context.Context, urls []*models.UrlDTO) ([]*models.UrlDTO, error) {
	batch := make([]*models.URL, 0, batchSize)
	shortUrls := make([]*models.UrlDTO, 0, len(urls))

	for _, url := range urls {
		shortUrl := tokengen.GenerateToken(tokenSize)

		batch = append(batch, &models.URL{
			CorrelationID: url.CorrelationID,
			BaseURL:       url.OriginURL,
			ShortURL:      shortUrl,
		})

		if len(batch) != batchSize {
			continue
		}

		dbUrls, err := uu.repo.BatchAddURL(ctx, batch)
		if err != nil {
			if len(shortUrls) == 0 {
				return nil, ErrNoBatchShorten
			}
			return shortUrls, ErrSomeBatchShortenFailed
		}
		batch = batch[:0]

		for _, u := range dbUrls {
			shortUrls = append(shortUrls, &models.UrlDTO{
				CorrelationID: u.CorrelationID,
				OriginURL:     u.BaseURL,
				ShortURL:      uu.getShortURL(u.ShortURL),
			})
		}
	}

	dbUrls, err := uu.repo.BatchAddURL(ctx, batch)
	if err != nil {
		if len(shortUrls) == 0 {
			return nil, ErrNoBatchShorten
		}
		return shortUrls, ErrSomeBatchShortenFailed
	}
	batch = batch[:0]

	for _, u := range dbUrls {
		shortUrls = append(shortUrls, &models.UrlDTO{
			CorrelationID: u.CorrelationID,
			OriginURL:     u.BaseURL,
			ShortURL:      uu.getShortURL(u.ShortURL),
		})
	}

	return shortUrls, nil
}

func (uu *UrlUsecase) GetURL(ctx context.Context, shortURL string) (string, bool) {
	url, err := uu.repo.GetURL(ctx, shortURL)
	if err != nil {
		uu.logger.Errorf("cannot get base_url for short_url=%s: %v", shortURL, err)
		return "", false
	}

	return url.BaseURL, true
}

func (uu *UrlUsecase) getShortURL(url string) string {
	return fmt.Sprintf("%s/%s", uu.cfg.BaseURL, url)
}
