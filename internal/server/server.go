package server

import (
	"net/http"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/MatiXxD/url-shortener/pkg/logger"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	mux    *chi.Mux
	cfg    *config.ServiceConfig
	logger *logger.Logger
}

func New(cfg *config.ServiceConfig, l *logger.Logger) *Server {
	return &Server{
		mux:    chi.NewRouter(),
		cfg:    cfg,
		logger: l,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Server running on %s\n", s.cfg.Addr)
	return http.ListenAndServe(s.cfg.Addr, s.mux)
}
