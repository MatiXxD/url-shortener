package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/MatiXxD/url-shortener/internal/repository"
)

func Router(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		ReduceURL(w, r)
	case http.MethodGet:
		GetURL(w, r)
	default:
		http.Error(w, "Wrong request", http.StatusBadRequest)
	}
}

func ReduceURL(w http.ResponseWriter, r *http.Request) {
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

	shortURL, err := repository.ReduceURL(string(url))
	if err != nil {
		http.Error(w, "Can't create short url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("http://localhost:8080/" + shortURL))
}

func GetURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.String()[1:]
	url, ok := repository.GetURL(shortURL)
	if !ok {
		http.Error(w, "Can't find url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
