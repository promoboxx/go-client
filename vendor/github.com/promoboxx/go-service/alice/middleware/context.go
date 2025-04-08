package middleware

import (
	"context"
	"database/sql"

	"github.com/promoboxx/go-auth/src/auth"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/sirupsen/logrus"
)

// getStringFromContext returns a string from the context or empty string if missing
func getStringFromContext(ctx context.Context, key contextkey.ContextKey) string {
	stri := ctx.Value(key)
	if str, ok := stri.(string); ok {
		return str
	}
	return ""
}

// GetRequestIDFromContext returns the requestID or empty string if missing
func GetRequestIDFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyRequestID)
}

// getInsecureUserIDFromContext returns the user id or empty string if missing
func GetInsecureUserIDFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyInsecureUserID)
}

// GetLoggerFromContext returns a logrus entry from the context, defaulting to a standard one if its not found
func GetLoggerFromContext(ctx context.Context) *logrus.Entry {
	inter := ctx.Value(contextkey.ContextKeyLogger)
	if entry, ok := inter.(*logrus.Entry); ok {
		return entry
	}
	return logrus.NewEntry(logrus.New()) // default to a new entry
}

func GetClaimsFromContext(ctx context.Context) auth.Claim {
	return ctx.Value(contextkey.ContextKeyClaims).(auth.Claim)
}

func GetDBConnFromContext(ctx context.Context) *sql.DB {
	return ctx.Value(contextkey.ContextKeyDBConn).(*sql.DB)
}

func GetDBFromContext[T any](ctx context.Context, f func(*sql.DB) T) T {
	if ctx.Value(contextkey.ContextKeyDB) == nil {
		return f(GetDBConnFromContext(ctx))
	}

	return ctx.Value(contextkey.ContextKeyDB).(T)
}

func GetCanaryVersionFromContext(ctx context.Context) string {
	return getStringFromContext(ctx, contextkey.ContextKeyCanary)
}
