package middleware

import (
	"context"
	"net/http"

	"github.com/promoboxx/go-service/alice/middleware/contextkey"
)

const (
	HeaderCanaryVersion = "x-canary-version"
)

func DetermineCanaryContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		version := r.Header.Get(HeaderCanaryVersion)
		if version != "" {
			ctx := context.WithValue(r.Context(), contextkey.ContextKeyCanary, version)
			r = r.WithContext(ctx)
		}

		h.ServeHTTP(w, r)
	})
}
