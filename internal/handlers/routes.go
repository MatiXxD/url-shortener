package handlers

import (
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/repository"
	"github.com/MatiXxD/url-shortener/internal/server"
	"github.com/MatiXxD/url-shortener/internal/usecase"
)

func BindRoutes(s *server.Server) {
	r := repository.NewMapRepository(map[string]*models.URL{})
	u := usecase.NewUrlUsecase(r)
	h := NewUrlHandler(u)

	s.Mux.HandleFunc("/", h.Router)
}
