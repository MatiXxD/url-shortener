package server

import "net/http"

type Server struct {
	Mux *http.ServeMux
}

func New() *Server {
	return &Server{http.NewServeMux()}
}

func (s *Server) Start() error {
	return http.ListenAndServe("localhost:8080", s.Mux)
}
