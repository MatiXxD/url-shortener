package main

import (
	"log"

	"github.com/MatiXxD/url-shortener/pkg/logger"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/MatiXxD/url-shortener/internal/server"
)

func main() {
	cfg := config.New()
	l, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(cfg, l)
	if err := s.BindRoutes(); err != nil {
		l.Fatalf("failed to bind routes: %v", err)
	}

	if err := s.Start(); err != nil {
		l.Fatalf("failed to start server: %v", err)
	}
}
