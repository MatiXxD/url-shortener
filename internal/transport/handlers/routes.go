package handlers

import (
	"github.com/MatiXxD/url-shortener/internal/server"
)

func BindRoutes(s *server.Server) {
	s.Mux.HandleFunc("/", Router)
}
