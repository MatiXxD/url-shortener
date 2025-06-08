package middleware

import (
	"net/http"
	"time"

	"github.com/MatiXxD/url-shortener/pkg/logger"
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

func LogMiddleware(logger *logger.Logger, h http.Handler) http.HandlerFunc {
	lf := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger = logger.With("request_id", GetRequestID(r.Context()))

		logger.Infof("got request: %s %s", r.Method, r.RequestURI)

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
		logger.Infof("request done: %s %s %d: size %d: time %s",
			r.Method,
			r.RequestURI,
			rd.status,
			rd.size,
			duration.String(),
		)
	}
	return lf
}
