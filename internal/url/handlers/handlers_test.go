package handlers

import (
	"bytes"
	"encoding/json"
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
		code     int
		response string
	}

	tests := []struct {
		name        string
		body        []byte
		contentType string
		want        want
	}{
		{
			name:        "Basic test",
			body:        []byte("/url"),
			contentType: "text/plain",
			want: want{
				code:     201,
				response: "http://localhost:8080/AAAAAAAA",
			},
		},
		{
			name:        "Wrong media type",
			body:        []byte("/"),
			contentType: "application/json",
			want: want{
				code:     415,
				response: "Wrong content type\n",
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
			require.Equal(t, tt.want.response, respBody)
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
				response: "",
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
				response: "Can't find url\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hdrs := []http.Header{}
			resp, respBody := createTestRequest(t, ts, http.MethodGet, tt.url, hdrs, bytes.NewBuffer(tt.body))
			require.Equal(t, tt.want.code, resp.StatusCode)
			require.Equal(t, tt.want.response, respBody)
			if !tt.isError {
				require.Equal(t, tt.want.location, resp.Header.Get("Location"))
			}
		})
	}
}

func TestUrlHandler_ShortenURL(t *testing.T) {
	d := map[string]*models.URL{}
	r := repository.NewMapRepository(d, l)
	mux, err := runTestServer(r)
	require.NoError(t, err)

	url := "/api/shorten"
	ts := httptest.NewServer(mux)

	type jsonResp struct {
		URL string `json:"url"`
	}

	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name        string
		body        []byte
		isError     bool
		contentType string
		want        want
	}{
		{
			name:        "Wrong content type",
			body:        []byte(`{"url": "https://google.com"}`),
			isError:     true,
			contentType: "wrong type",
			want: want{
				code:     415,
				response: "Wrong content type\n",
			},
		},
		{
			name:        "Empty body",
			body:        []byte(``),
			isError:     true,
			contentType: "application/json",
			want: want{
				code:     500,
				response: "Can't read body\n",
			},
		},
		{
			name:        "ShortenURL OK",
			body:        []byte(`{"url": "https://google.com"}`),
			isError:     false,
			contentType: "application/json",
			want: want{
				code: 200,
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
			resp, respBody := createTestRequest(t, ts, http.MethodPost, url, hdrs, bytes.NewBuffer(tt.body))
			require.Equal(t, tt.want.code, resp.StatusCode)
			if tt.isError {
				require.Equal(t, tt.want.response, respBody)
			} else {
				jsonBody := &jsonResp{}
				err := json.Unmarshal([]byte(respBody), jsonBody)
				require.NoError(t, err)
				require.True(t, len(jsonBody.URL) > 0)
				_, ok := r.GetURL(jsonBody.URL)
				require.True(t, ok)
			}
		})
	}
}
