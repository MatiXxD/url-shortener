package server

import (
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/handlers"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
)

func BindRoutes(s *Server) {
	r := repository.NewMapRepository(map[string]*models.URL{})
	u := usecase.NewUrlUsecase(r)
	h := handlers.NewUrlHandler(u, s.Cfg)

	s.Mux.Post("/", h.ReduceURL)
	s.Mux.Get("/{url}", h.GetURL)
}
