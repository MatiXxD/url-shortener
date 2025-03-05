package main

import (
	"log"

	"github.com/MatiXxD/url-shortener/config"
	"github.com/MatiXxD/url-shortener/internal/server"
)

func main() {
	cfg := config.New()
	s := server.New(cfg)
	server.BindRoutes(s)
	if err := s.Start(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
