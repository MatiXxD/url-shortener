package main

import (
	"log"

	"github.com/MatiXxD/url-shortener/internal/handlers"
	"github.com/MatiXxD/url-shortener/internal/server"
)

func main() {
	s := server.New()
	handlers.BindRoutes(s)
	if err := s.Start(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}
