package repository

import (
	"context"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestMapRepository_AddURL(t *testing.T) {
	testURL := "https://www.google.com"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: {
			BaseURL:  testURL,
			ShortURL: testShortURL,
		},
	}
	repo := NewMapRepository(d, l)

	t.Run("Success add", func(t *testing.T) {
		url := "https://ya.ru"
		shortURL := "BBBBB"
		corrID := uuid.New().String()
		got, err := repo.AddURL(context.Background(), &models.URL{
			CorrelationID: corrID,
			BaseURL:       url,
			ShortURL:      shortURL,
		})

		require.NoError(t, err)
		require.Equal(t, shortURL, got)
	})

	t.Run("Alredy exists", func(t *testing.T) {
		shortURL := "CCCCC"
		corrID := uuid.New().String()
		got, err := repo.AddURL(context.Background(), &models.URL{
			CorrelationID: corrID,
			BaseURL:       testURL,
			ShortURL:      shortURL,
		})

		require.NoError(t, err)
		require.Equal(t, testShortURL, got)
	})
}

func TestMapRepository_GetURL(t *testing.T) {
	testURL := "https://www.google.com"
	testShortURL := "AAAAA"
	corrID := uuid.New().String()
	d := map[string]*models.URL{
		testURL: (&models.URL{
			CorrelationID: corrID,
			BaseURL:       testURL,
			ShortURL:      testShortURL,
		}),
	}
	repo := NewMapRepository(d, l)

	t.Run("Success get", func(t *testing.T) {
		getURL, err := repo.GetURL(context.Background(), testShortURL)
		require.NoError(t, err)
		require.Equal(t, testURL, getURL.BaseURL)
	})

	t.Run("Can't get url", func(t *testing.T) {
		shortURL := "https://www.random.com"
		_, err := repo.GetURL(context.Background(), shortURL)
		require.ErrorContains(t, err, "not found")
	})
}
