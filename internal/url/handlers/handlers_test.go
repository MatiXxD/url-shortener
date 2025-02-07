package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
	"github.com/stretchr/testify/require"
)

func TestUrlHandler_ReduceURL(t *testing.T) {
	d := map[string]*models.URL{
		"/url": {BaseURL: "/url", ShortURL: "AAAAAAAA"},
	}
	r := repository.NewMapRepository(d)
	u := usecase.NewUrlUsecase(r)
	h := NewUrlHandler(u)

	type want struct {
		code        int
		contentType string
		response    string
	}
	tests := []struct {
		name        string
		body        string
		contentType string
		isError     bool
		want        want
	}{
		{
			name:        "Basic test",
			body:        "/url",
			contentType: "text/plain",
			isError:     false,
			want: want{
				code:        201,
				contentType: "text/plain",
				response:    "http://localhost:8080/AAAAAAAA",
			},
		},
		{
			name:        "Wrong media type",
			body:        "/",
			contentType: "application/json",
			isError:     true,
			want: want{
				code:        415,
				contentType: "text/plain",
				response:    "Wrong content type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			r.Header.Add("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			h.ReduceURL(w, r)

			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode)

			if !tt.isError {
				defer res.Body.Close()
				resBody, err := io.ReadAll(res.Body)

				require.NoError(t, err)
				require.Equal(t, tt.want.response, string(resBody))
				require.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestUrlHandler_GetURL(t *testing.T) {
	d := map[string]*models.URL{
		"/url": {BaseURL: "/url", ShortURL: "AAAAAAAA"},
	}
	r := repository.NewMapRepository(d)
	u := usecase.NewUrlUsecase(r)
	h := NewUrlHandler(u)

	type want struct {
		code     int
		location string
		response string
	}
	tests := []struct {
		name    string
		url     string
		isError bool
		want    want
	}{
		{
			name:    "Basic test",
			url:     "/AAAAAAAA",
			isError: false,
			want: want{
				code:     307,
				location: "/url",
				response: "http://localhost:8080/AAAAAAAA",
			},
		},
		{
			name:    "Can't find url",
			url:     "/random",
			isError: true,
			want: want{
				code:     400,
				location: "",
				response: "Wrong content type",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			h.GetURL(w, r)

			res := w.Result()
			require.Equal(t, tt.want.code, res.StatusCode)

			if !tt.isError {
				require.Equal(t, tt.want.location, res.Header.Get("Location"))
			}
		})
	}
}
