package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/MatiXxD/url-shortener/config"
	mw "github.com/MatiXxD/url-shortener/internal/middleware"
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

	shortURL, err := uh.urlUsecase.ReduceURL(r.Context(), &models.UrlDTO{
		CorrelationID: uuid.New().String(),
		OriginURL:     string(url),
	})
	if err != nil {
		logger.Error("can't create short URL")
		http.Error(w, "Can't create short url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

func (uh *UrlHandler) BatchReduceURL(w http.ResponseWriter, r *http.Request) {
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

	var urls []*models.UrlDTO
	if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
		logger.Error("can't unmarshal request body")
		http.Error(w, "Can't read body", http.StatusInternalServerError)
		return
	}

	shortUrls, err := uh.urlUsecase.BatchReduceURL(r.Context(), urls)
	if err != nil {
		logger.Error("can't short all urls")
		http.Error(w, "Can't create short urls", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(shortUrls); err != nil {
		logger.Error("can't marshal response body")
		http.Error(w, "Can't marshal response body", http.StatusInternalServerError)
		return
	}
}

func (uh *UrlHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	logger := uh.logger
	reqID := mw.GetRequestID(r.Context())
	if reqID != "" {
		logger = uh.logger.With("request_id", reqID)
	}

	shortURL := chi.URLParam(r, "url")
	url, ok := uh.urlUsecase.GetURL(r.Context(), shortURL)
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

	var reqUrl models.ShortenURLReqBody
	if err := easyjson.UnmarshalFromReader(r.Body, &reqUrl); err != nil {
		logger.Error("can't unmarshal request body")
		http.Error(w, "Can't read body", http.StatusInternalServerError)
		return
	}

	shortUrl, err := uh.urlUsecase.ReduceURL(r.Context(), &models.UrlDTO{
		CorrelationID: uuid.New().String(),
		OriginURL:     reqUrl.URL,
	})
	if err != nil {
		logger.Error("can't create short URL")
		http.Error(w, "Can't create short url", http.StatusInternalServerError)
		return
	}

	resp := &models.ShortenURLRespBody{
		ShortURL: shortUrl,
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := easyjson.MarshalToWriter(resp, w); err != nil {
		logger.Error("can't marshal response body")
		http.Error(w, "Can't marshal response body", http.StatusInternalServerError)
		return
	}
}
