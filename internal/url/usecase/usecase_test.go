package usecase

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"go.uber.org/zap"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/stretchr/testify/require"
)

var (
	l   *logger.Logger
	cfg *config.ServiceConfig
)

func TestMain(t *testing.M) {
	zl, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	l = &logger.Logger{zl.Sugar()}

	cfg = &config.ServiceConfig{
		Addr:        ":8080",
		BaseURL:     "http://localhost:8080",
		LoggerLevel: "info",
		FilePath:    "/tmp/short-url-db.json",
		DSN:         "",
	}

	os.Exit(t.Run())
}

func TestUsecase_ReduceURL(t *testing.T) {
	testURL := "https://www.ya.ru"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: {
			BaseURL:  testURL,
			ShortURL: testShortURL,
		},
	}
	r := repository.NewMapRepository(d, l)
	uc := NewUrlUsecase(r, cfg, l)

	t.Run("Success add", func(t *testing.T) {
		url := "https://www.google.com"
		shortURL, err := uc.ReduceURL(context.Background(), &models.UrlDTO{
			OriginURL: url,
		})

		require.NoError(t, err)
		require.NotZero(t, shortURL)
	})

	t.Run("Alredy exists", func(t *testing.T) {
		shortURL, err := uc.ReduceURL(context.Background(), &models.UrlDTO{
			OriginURL: testURL,
		})

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s/%s", cfg.BaseURL, testShortURL), shortURL)
	})
}

func TestUsecase_GetURL(t *testing.T) {
	testURL := "https://www.google.com"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: {
			BaseURL:  testURL,
			ShortURL: testShortURL,
		},
	}
	r := repository.NewMapRepository(d, l)
	uc := NewUrlUsecase(r, cfg, l)

	t.Run("Success get", func(t *testing.T) {
		got, ok := uc.GetURL(context.Background(), testShortURL)
		require.Equal(t, true, ok)
		require.Equal(t, testURL, got)
	})

	t.Run("Can't get url", func(t *testing.T) {
		got, ok := uc.GetURL(context.Background(), "https://random.com")
		require.Equal(t, false, ok)
		require.Zero(t, got)
	})
}
