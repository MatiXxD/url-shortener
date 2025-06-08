package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKeyRequestID struct{}

func RequestIdMiddleware(h http.Handler) http.Handler {
	mw := func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.New()
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, reqID)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(mw)
}

func GetRequestID(ctx context.Context) string {
	rawReqID := ctx.Value(ctxKeyRequestID{})

	reqID, ok := rawReqID.(uuid.UUID)
	if !ok {
		return ""
	}

	return reqID.String()
}
