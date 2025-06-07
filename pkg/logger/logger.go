package logger

import (
	"github.com/MatiXxD/url-shortener/config"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(serviceCfg *config.ServiceConfig) (*Logger, error) {
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

	return &Logger{logger.Sugar()}, nil
}
