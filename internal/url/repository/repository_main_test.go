package repository

import (
	"os"
	"testing"

	"github.com/MatiXxD/url-shortener/pkg/logger"
	"go.uber.org/zap"
)

var l *logger.Logger

func TestMain(t *testing.M) {
	zl, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	l = &logger.Logger{zl.Sugar()}

	os.Exit(t.Run())
}
