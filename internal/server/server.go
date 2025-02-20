package server

import (
	"log"
	"net/http"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	Mux *chi.Mux
	Cfg *config.ServiceConfig
}

func New(cfg *config.ServiceConfig) *Server {
	return &Server{
		Mux: chi.NewRouter(),
		Cfg: cfg,
	}
}

func (s *Server) Start() error {
	log.Printf("Server running on %s\n", s.Cfg.Addr)
	return http.ListenAndServe(s.Cfg.Addr, s.Mux)
}
