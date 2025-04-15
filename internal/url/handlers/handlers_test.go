package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/stretchr/testify/require"
)

func TestUrlHandler_ReduceURL(t *testing.T) {
	d := map[string]*models.URL{
		"/url": {BaseURL: "/url", ShortURL: "AAAAAAAA"},
	}
	r := repository.NewMapRepository(d, l)
	mux, err := runTestServer(r)
	require.NoError(t, err)

	ts := httptest.NewServer(mux)

	type want struct {
		code        int
		contentType string
		response    string
	}

	tests := []struct {
		name        string
		body        []byte
		contentType string
		isError     bool
		want        want
	}{
		{
			name:        "Basic test",
			body:        []byte("/url"),
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
			body:        []byte("/"),
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
			hdrs := []http.Header{
				{
					"Content-Type": []string{tt.contentType},
				},
			}
			resp, respBody := createTestRequest(t, ts, http.MethodPost, "/", hdrs, bytes.NewBuffer(tt.body))
			require.Equal(t, tt.want.code, resp.StatusCode)
			if !tt.isError {
				require.NoError(t, err)
				require.Equal(t, tt.want.response, respBody)
				require.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))
			}
		})
	}
}

func TestUrlHandler_GetURL(t *testing.T) {
	d := map[string]*models.URL{
		"/url": {BaseURL: "/url", ShortURL: "AAAAAAAA"},
	}
	r := repository.NewMapRepository(d, l)
	mux, err := runTestServer(r)
	require.NoError(t, err)

	ts := httptest.NewServer(mux)

	type want struct {
		code     int
		location string
		response string
	}
	tests := []struct {
		name    string
		body    []byte
		url     string
		isError bool
		want    want
	}{
		{
			name:    "Basic test",
			body:    []byte(""),
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
			body:    []byte(""),
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
			hdrs := []http.Header{}
			resp, _ := createTestRequest(t, ts, http.MethodGet, tt.url, hdrs, bytes.NewBuffer(tt.body))
			require.Equal(t, tt.want.code, resp.StatusCode)
			if !tt.isError {
				require.Equal(t, tt.want.location, resp.Header.Get("Location"))
			}
		})
	}
}
