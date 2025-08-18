package repository

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFileRepository_initCache(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		wantCache   map[string]*models.URL
		wantErr     bool
	}{
		{
			name:        "valid JSON",
			fileContent: `{"id": 1, "correlation_id":"5f2ee353-1946-4b78-8801-eebd5ca17aee","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}`,
			wantCache:   map[string]*models.URL{"http://yandex.ru": {ID: 1, CorrelationID: "5f2ee353-1946-4b78-8801-eebd5ca17aee", ShortURL: "4rSPg8ap", BaseURL: "http://yandex.ru"}},
			wantErr:     false,
		},
		{
			name: "multiple JSON",
			fileContent: `{"id": 1, "correlation_id":"5f2ee353-1946-4b78-8801-eebd5ca17aee","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}
			{"id": 2,  "correlation_id":"3d17d9d7-68a4-40b4-a525-47c64383f3e3","short_url":"edVPg3ks","original_url":"http://ya.ru"}`,
			wantCache: map[string]*models.URL{
				"http://yandex.ru": {ID: 1, CorrelationID: "5f2ee353-1946-4b78-8801-eebd5ca17aee", ShortURL: "4rSPg8ap", BaseURL: "http://yandex.ru"},
				"http://ya.ru":     {ID: 2, CorrelationID: "3d17d9d7-68a4-40b4-a525-47c64383f3e3", ShortURL: "edVPg3ks", BaseURL: "http://ya.ru"},
			},
			wantErr: false,
		},
		{
			name:        "invalid JSON",
			fileContent: `{"id":}`,
			wantCache:   nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_cache_*.json")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			if tt.fileContent != "" {
				_, err = tmpFile.WriteString(tt.fileContent)
				require.NoError(t, err)
				require.NoError(t, tmpFile.Sync())
			}
			require.NoError(t, tmpFile.Close())

			fr, err := NewFileRepository(tmpFile.Name(), l)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				cacheJSON, _ := json.Marshal(fr.cache)
				wantJSON, _ := json.Marshal(tt.wantCache)
				require.JSONEq(t, string(wantJSON), string(cacheJSON))
			}
		})
	}
}

func TestFileRepository_saveURL(t *testing.T) {
	type testCase struct {
		name     string
		input    *models.URL
		expected *models.URL
	}

	tests := []testCase{
		{
			name: "simple valid URL",
			input: &models.URL{
				ID:            1,
				CorrelationID: "5f2ee353-1946-4b78-8801-eebd5ca17aee",
				ShortURL:      "abc123",
				BaseURL:       "http://example.com",
			},
			expected: &models.URL{
				ID:            1,
				CorrelationID: "5f2ee353-1946-4b78-8801-eebd5ca17aee",
				ShortURL:      "abc123",
				BaseURL:       "http://example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "test_save_url_*.json")
			require.NoError(t, err)
			defer os.Remove(tmpFile.Name())

			fr, err := NewFileRepository(tmpFile.Name(), l)

			err = fr.saveURL(tt.input)
			require.NoError(t, err)

			data, err := os.ReadFile(tmpFile.Name())
			require.NoError(t, err)

			var result models.URL
			err = json.Unmarshal(data, &result)
			require.NoError(t, err)

			require.Equal(t, tt.expected.ID, result.ID)
			require.Equal(t, tt.expected.ShortURL, result.ShortURL)
			require.Equal(t, tt.expected.BaseURL, result.BaseURL)
		})
	}
}

func TestFileRepository_AddURL(t *testing.T) {
	type testCase struct {
		name          string
		initialCache  map[string]*models.URL
		inputURL      string
		inputShortURL string
		wantShortURL  string
		expectNew     bool
	}

	tests := []testCase{
		{
			name:          "new URL added",
			initialCache:  map[string]*models.URL{},
			inputURL:      "http://example.com",
			inputShortURL: "abc123",
			wantShortURL:  "abc123",
			expectNew:     true,
		},
		{
			name: "URL already in cache",
			initialCache: map[string]*models.URL{
				"http://example.com": {
					ID:            1,
					CorrelationID: uuid.New().String(),
					BaseURL:       "http://example.com",
					ShortURL:      "cached123",
				},
			},
			inputURL:      "http://example.com",
			inputShortURL: "ignored123",
			wantShortURL:  "cached123",
			expectNew:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr, err := NewFileRepository("", l)
			require.NoError(t, err)

			fr.cache = tt.initialCache

			got, err := fr.AddURL(context.Background(), &models.URL{
				BaseURL:  tt.inputURL,
				ShortURL: tt.inputShortURL,
			})

			require.NoError(t, err)
			require.Equal(t, tt.wantShortURL, got)

			u, ok := fr.cache[tt.inputURL]
			require.True(t, ok)

			if tt.expectNew {
				require.Equal(t, tt.inputShortURL, u.ShortURL)
				require.Equal(t, tt.inputURL, u.BaseURL)
			} else {
				require.Equal(t, "cached123", u.ShortURL)
			}
		})
	}
}

func TestFileRepository_GetURL(t *testing.T) {
	type testCase struct {
		name         string
		cache        map[string]*models.URL
		inputShort   string
		wantOriginal string
		wantFound    bool
	}

	tests := []testCase{
		{
			name: "short URL exists",
			cache: map[string]*models.URL{
				"http://example.com": {
					ID:            1,
					CorrelationID: uuid.New().String(),
					BaseURL:       "http://example.com",
					ShortURL:      "abc123",
				},
			},
			inputShort:   "abc123",
			wantOriginal: "http://example.com",
			wantFound:    true,
		},
		{
			name: "short URL does not exist",
			cache: map[string]*models.URL{
				"http://example.com": {
					ID:            1,
					CorrelationID: uuid.New().String(),
					BaseURL:       "http://example.com",
					ShortURL:      "abc123",
				},
			},
			inputShort:   "not_found",
			wantOriginal: "",
			wantFound:    false,
		},
		{
			name:         "empty cache",
			cache:        map[string]*models.URL{},
			inputShort:   "abc123",
			wantOriginal: "",
			wantFound:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr, err := NewFileRepository("", l)
			require.NoError(t, err)

			fr.cache = tt.cache

			gotURL, err := fr.GetURL(context.Background(), tt.inputShort)
			if !tt.wantFound {
				require.ErrorContains(t, err, "not found")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantOriginal, gotURL.BaseURL)

			}
		})
	}
}
