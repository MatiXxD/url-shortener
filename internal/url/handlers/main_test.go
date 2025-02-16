package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func createTestRequest(t *testing.T, ts *httptest.Server,
	method, path string, headers []http.Header, body io.Reader,
) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	for _, header := range headers {
		for k, v := range header {
			for _, vv := range v {
				req.Header.Set(k, vv)
			}
		}
	}

	// using last response for tests
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func runTestServer(r url.Repository) (chi.Router, error) {
	u := usecase.NewUrlUsecase(r)
	h := NewUrlHandler(u)

	mux := chi.NewRouter()
	mux.Post("/", h.ReduceURL)
	mux.Get("/{url}", h.GetURL)

	return mux, nil
}
