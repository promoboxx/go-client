package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
)

const (
	HeaderRequestID = "x-request-id"
)

// RequestID adds a request id to the context, if available it will use the one in the x-request-id header
func RequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//check header for existing request id
		rID := r.Header.Get(HeaderRequestID)
		if len(rID) == 0 {
			//generate id
			rID = uuid.New().String()
		}

		// add rID to the context
		ctx := context.WithValue(r.Context(), contextkey.ContextKeyRequestID, rID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
