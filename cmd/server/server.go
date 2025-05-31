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
	s.BindRoutes()
	if err := s.Start(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
