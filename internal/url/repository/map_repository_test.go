package repository

import (
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/stretchr/testify/require"
)

func TestMapRepository_ReduceURL(t *testing.T) {
	d := map[string]*models.URL{}
	repo := NewMapRepository(d)

	t.Run("Success add", func(t *testing.T) {
		url := "https://www.google.com"
		shortURL, err := repo.ReduceURL(url)

		require.NoError(t, err)
		require.NotZero(t, shortURL)
	})
}

func TestMapRepository_GetURL(t *testing.T) {
	d := map[string]*models.URL{}
	repo := NewMapRepository(d)

	t.Run("Success get", func(t *testing.T) {
		url := "https://www.google.com"
		shortURL, _ := repo.ReduceURL(url)

		getURL, ok := repo.GetURL(shortURL)

		require.Equal(t, true, ok)
		require.Equal(t, url, getURL)
	})

	t.Run("Can't get url", func(t *testing.T) {
		shortURL := "https://www.random.com"

		getURL, ok := repo.GetURL(shortURL)

		require.Equal(t, false, ok)
		require.Zero(t, getURL)
	})
}
