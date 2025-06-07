package server

import (
	"net/http"

	mw "github.com/MatiXxD/url-shortener/internal/middleware"
	"github.com/MatiXxD/url-shortener/internal/url/handlers"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
)

type middleware func(http.Handler) http.Handler

func (s *Server) BindRoutes() error {
	r, err := repository.NewFileRepository(s.cfg.FilePath, s.logger)
	if err != nil {
		s.logger.Errorf("failed to create repository: %v", err)
		return err
	}

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

	return nil
}
