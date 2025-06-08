package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/MatiXxD/url-shortener/config"
	mw "github.com/MatiXxD/url-shortener/internal/middleware"
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
)

type UrlHandler struct {
	urlUsecase url.Usecase
	cfg        *config.ServiceConfig
	logger     *logger.Logger
}

func NewUrlHandler(u url.Usecase, cfg *config.ServiceConfig, l *logger.Logger) *UrlHandler {
	return &UrlHandler{
		urlUsecase: u,
		cfg:        cfg,
		logger:     l,
	}
}

func (uh *UrlHandler) ReduceURL(w http.ResponseWriter, r *http.Request) {
	logger := uh.logger
	reqID := mw.GetRequestID(r.Context())
	if reqID != "" {
		logger = uh.logger.With("request_id", reqID)
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/plain") {
		logger.Error("request contains wrong content type")
		http.Error(w, "Wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	url, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error("request doesn't have body")
		http.Error(w, "Can't read request body", http.StatusBadRequest)
		return
	}

	shortURL, err := uh.urlUsecase.ReduceURL(string(url))
	if err != nil {
		logger.Error("can't create short URL")
		http.Error(w, "Can't create short url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(uh.cfg.BaseURL + "/" + shortURL))
}

func (uh *UrlHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	logger := uh.logger
	reqID := mw.GetRequestID(r.Context())
	if reqID != "" {
		logger = uh.logger.With("request_id", reqID)
	}

	shortURL := chi.URLParam(r, "url")
	url, ok := uh.urlUsecase.GetURL(shortURL)
	if !ok {
		logger.Error("can't find url")
		http.Error(w, "Can't find url", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (uh *UrlHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	logger := uh.logger
	reqID := mw.GetRequestID(r.Context())
	if reqID != "" {
		logger = uh.logger.With("request_id", reqID)
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		logger.Error("request contains wrong content type")
		http.Error(w, "Wrong content type", http.StatusUnsupportedMediaType)
		return
	}

	reqUrl := new(models.UrlDTO)
	if err := easyjson.UnmarshalFromReader(r.Body, reqUrl); err != nil {
		logger.Error("can't unmarshal request body")
		http.Error(w, "Can't read body", http.StatusInternalServerError)
		return
	}

	newUrl, err := uh.urlUsecase.ReduceURL(reqUrl.URL)
	if err != nil {
		logger.Error("can't create short URL")
		http.Error(w, "Can't create short url", http.StatusInternalServerError)
		return
	}
	newUrlDTO := models.NewUrlDTO(newUrl)

	w.Header().Set("Content-Type", "application/json")
	if _, err := easyjson.MarshalToWriter(newUrlDTO, w); err != nil {
		logger.Error("can't marshal response body")
		http.Error(w, "Can't marshal response body", http.StatusInternalServerError)
		return
	}
}
