package repository

import (
	"go.uber.org/zap"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/stretchr/testify/require"
)

var l *zap.Logger

func TestMain(t *testing.M) {
	var err error
	l, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

func TestMapRepository_AddURL(t *testing.T) {
	testURL := "https://www.google.com"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: models.NewURL(testURL, testShortURL),
	}
	repo := NewMapRepository(d, l)

	t.Run("Success add", func(t *testing.T) {
		url := "https://ya.ru"
		shortURL := "BBBBB"
		got, err := repo.AddURL(url, shortURL)
		require.NoError(t, err)
		require.Equal(t, shortURL, got)
	})

	t.Run("Alredy exists", func(t *testing.T) {
		shortURL := "CCCCC"
		got, err := repo.AddURL(testURL, shortURL)
		require.NoError(t, err)
		require.Equal(t, testShortURL, got)
	})
}

func TestMapRepository_GetURL(t *testing.T) {
	testURL := "https://www.google.com"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: models.NewURL(testURL, testShortURL),
	}
	repo := NewMapRepository(d, l)

	t.Run("Success get", func(t *testing.T) {
		getURL, ok := repo.GetURL(testShortURL)
		require.Equal(t, true, ok)
		require.Equal(t, testURL, getURL)
	})

	t.Run("Can't get url", func(t *testing.T) {
		shortURL := "https://www.random.com"
		getURL, ok := repo.GetURL(shortURL)
		require.Equal(t, false, ok)
		require.Zero(t, getURL)
	})
}
