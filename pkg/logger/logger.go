package logger

import (
	"github.com/MatiXxD/url-shortener/config"
	"go.uber.org/zap"
)

func NewLogger(serviceCfg *config.ServiceConfig) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(serviceCfg.LoggerLevel)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
