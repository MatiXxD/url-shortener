package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/go-chi/chi/v5"
)

type UrlHandler struct {
	urlUsecase url.Usecase
}

func NewUrlHandler(u url.Usecase) *UrlHandler {
	return &UrlHandler{
		urlUsecase: u,
	}
}

func (uh *UrlHandler) ReduceURL(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		http.Error(w, "Wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	url, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Can't read request body", http.StatusBadRequest)
		return
	}

	shortURL, err := uh.urlUsecase.ReduceURL(string(url))
	if err != nil {
		http.Error(w, "Can't create short url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("http://localhost:8080/" + shortURL))
}

func (uh *UrlHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	shortURL := chi.URLParam(r, "url")
	url, ok := uh.urlUsecase.GetURL(shortURL)
	if !ok {
		http.Error(w, "Can't find url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
