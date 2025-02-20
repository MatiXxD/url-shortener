package config

import "flag"

type ServiceConfig struct {
	Addr    string
	BaseURL string
}

func ParseConfig() *ServiceConfig {
	cfg := &ServiceConfig{}

	flag.StringVar(&cfg.Addr, "a", ":8080", "Addres and port for server")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "BaseURL for short ulrs")

	flag.Parse()

	return cfg
}
