package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.data.size += size
	return size, err
}

func (w *loggingResponseWriter) WriteHeader(status int) {
	w.ResponseWriter.WriteHeader(status)
	w.data.status = status
}

func LogMiddleware(logger *zap.Logger, h http.Handler) http.HandlerFunc {
	lf := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rd := &responseData{
			status: http.StatusOK,
			size:   0,
		}
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			data:           rd,
		}
		h.ServeHTTP(lw, r)

		duration := time.Since(start)
		logger.Info(fmt.Sprintf("[%s] %s %d: size %d: time %s",
			r.Method,
			r.RequestURI,
			rd.status,
			rd.size,
			duration.String()),
		)
	}
	return lf
}
