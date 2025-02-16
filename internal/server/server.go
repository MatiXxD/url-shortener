package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	Mux *chi.Mux
}

func New() *Server {
	return &Server{chi.NewRouter()}
}

func (s *Server) Start() error {
	return http.ListenAndServe("localhost:8080", s.Mux)
}
