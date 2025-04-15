package usecase

import (
	"go.uber.org/zap"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
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

func TestUsecase_ReduceURL(t *testing.T) {
	testURL := "https://www.ya.ru"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: models.NewURL(testURL, testShortURL),
	}
	r := repository.NewMapRepository(d, l)
	uc := NewUrlUsecase(r, l)

	t.Run("Success add", func(t *testing.T) {
		url := "https://www.google.com"
		shortURL, err := uc.ReduceURL(url)
		require.NoError(t, err)
		require.NotZero(t, shortURL)
	})

	t.Run("Alredy exists", func(t *testing.T) {
		shortURL, err := uc.ReduceURL(testURL)
		require.NoError(t, err)
		require.Equal(t, testShortURL, shortURL)
	})
}

func TestUsecase_GetURL(t *testing.T) {
	testURL := "https://www.google.com"
	testShortURL := "AAAAA"
	d := map[string]*models.URL{
		testURL: models.NewURL(testURL, testShortURL),
	}
	r := repository.NewMapRepository(d, l)
	uc := NewUrlUsecase(r, l)

	t.Run("Success get", func(t *testing.T) {
		got, ok := uc.GetURL(testShortURL)
		require.Equal(t, true, ok)
		require.Equal(t, testURL, got)
	})

	t.Run("Can't get url", func(t *testing.T) {
		got, ok := uc.GetURL("https://random.com")
		require.Equal(t, false, ok)
		require.Zero(t, got)
	})
}
