package server

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	mux    *chi.Mux
	cfg    *config.ServiceConfig
	logger *zap.Logger
}

func New(cfg *config.ServiceConfig, l *zap.Logger) *Server {
	return &Server{
		mux:    chi.NewRouter(),
		cfg:    cfg,
		logger: l,
	}
}

func (s *Server) Start() error {
	log.Printf("Server running on %s\n", s.cfg.Addr)
	return http.ListenAndServe(s.cfg.Addr, s.mux)
}
