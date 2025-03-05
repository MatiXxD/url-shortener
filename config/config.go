package config

import (
	"flag"
	"os"
)

type ServiceConfig struct {
	Addr    string
	BaseURL string
}

const (
	defaultAddr    = ":8080"
	defaultBaseURL = "http://localhost:8080"
)

func New() *ServiceConfig {
	cfg := &ServiceConfig{}

	flag.StringVar(&cfg.Addr, "a", defaultAddr, "Addres and port for server")
	flag.StringVar(&cfg.BaseURL, "b", defaultBaseURL, "BaseURL for short ulrs")
	flag.Parse()

	parseEnv(cfg)

	return cfg
}

func parseEnv(cfg *ServiceConfig) {
	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.Addr = addr
	}
	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.Addr = baseURL
	}
}
