package server

import (
	"net/http"

	mw "github.com/MatiXxD/url-shortener/internal/middleware"
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/handlers"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
)

type middleware func(http.Handler) http.Handler

func (s *Server) BindRoutes() {
	r := repository.NewMapRepository(map[string]*models.URL{}, s.logger)
	u := usecase.NewUrlUsecase(r, s.logger)
	h := handlers.NewUrlHandler(u, s.cfg, s.logger)

	logMiddleware := func(next http.Handler) http.Handler {
		return mw.LogMiddleware(s.logger, next)
	}

	middlewares := []middleware{
		logMiddleware,
		mw.CompressMiddleware,
	}

	for _, m := range middlewares {
		s.mux.Use(m)
	}

	s.mux.Post("/", h.ReduceURL)
	s.mux.Get("/{url}", h.GetURL)

	s.mux.Post("/api/shorten", h.ShortenURL)
}
