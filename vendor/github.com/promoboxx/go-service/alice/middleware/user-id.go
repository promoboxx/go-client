package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/promoboxx/go-service/alice/middleware/contextkey"
)

type JWTDecoder interface {
	DecodeSegment(seg string) ([]byte, error)
}

// UserIDIngector injects a UserID from a JWT into the context
type UserIDIngector interface {
	Inject(h http.Handler) http.Handler
}

type userIDInjector struct {
	jwtDecoder JWTDecoder
}

// NewLogrusLogger allows you to setup a base entry to use for logging
func NewUserIDInjector(jwt JWTDecoder) UserIDIngector {
	return &userIDInjector{jwtDecoder: jwt}
}

func (u *userIDInjector) Inject(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID string
		t := r.Header.Get("Authorization")
		t = strings.TrimPrefix(t, "Bearer ")
		parts := strings.Split(t, ".")

		if len(parts) == 3 {
			by, err := u.jwtDecoder.DecodeSegment(parts[1])
			if err == nil {
				claims := make(map[string]interface{})
				err = json.Unmarshal(by, &claims)
				if err == nil {
					if userIDi, ok := claims["sub"]; ok {
						if v, ok := userIDi.(string); ok {
							userID = v
						}
					}

				}
			}
		}

		// add userID to the context
		ctx := context.WithValue(r.Context(), contextkey.ContextKeyInsecureUserID, userID)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
