package server

import (
	"net/http"

	mw "github.com/MatiXxD/url-shortener/internal/middleware"
	"github.com/MatiXxD/url-shortener/internal/url"
	"github.com/MatiXxD/url-shortener/internal/url/handlers"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
	"github.com/MatiXxD/url-shortener/pkg/postgres"
)

type middleware func(http.Handler) http.Handler

func (s *Server) BindRoutes() error {
	var (
		err error
		r   url.Repository
	)

	if s.cfg.DSN != "" {
		db, err := postgres.New(s.cfg.DSN)
		if err != nil {
			s.logger.Errorf("failed to connect to postgres: %v", err)
			return err
		}

		r = repository.NewPostgresRepository(db, s.logger)
	} else {
		r, err = repository.NewFileRepository(s.cfg.FilePath, s.logger)
		if err != nil {
			s.logger.Errorf("failed to create repository: %v", err)
			return err
		}
	}

	u := usecase.NewUrlUsecase(r, s.cfg, s.logger)
	h := handlers.NewUrlHandler(u, s.cfg, s.logger)

	logMiddleware := func(next http.Handler) http.Handler {
		return mw.LogMiddleware(s.logger, next)
	}

	middlewares := []middleware{
		mw.RequestIdMiddleware,
		logMiddleware,
		mw.CompressMiddleware,
	}

	for _, m := range middlewares {
		s.mux.Use(m)
	}

	s.mux.Post("/", h.ReduceURL)
	s.mux.Get("/{url}", h.GetURL)
	s.mux.Post("/api/shorten", h.ShortenURL)
	s.mux.Post("/api/shorten/batch", h.BatchReduceURL)

	return nil
}
