package server

import (
	"github.com/MatiXxD/url-shortener/internal/middleware"
	"github.com/MatiXxD/url-shortener/internal/models"
	"github.com/MatiXxD/url-shortener/internal/url/handlers"
	"github.com/MatiXxD/url-shortener/internal/url/repository"
	"github.com/MatiXxD/url-shortener/internal/url/usecase"
	"net/http"
)

func (s *Server) BindRoutes() {
	r := repository.NewMapRepository(map[string]*models.URL{}, s.logger)
	u := usecase.NewUrlUsecase(r, s.logger)
	h := handlers.NewUrlHandler(u, s.cfg, s.logger)

	s.mux.Post("/", middleware.LogMiddleware(s.logger, http.HandlerFunc(h.ReduceURL)))
	s.mux.Get("/{url}", middleware.LogMiddleware(s.logger, http.HandlerFunc(h.GetURL)))
}
