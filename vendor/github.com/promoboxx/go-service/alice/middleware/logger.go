package middleware

import (
	"context"
	"net/http"

	"github.com/promoboxx/go-service/alice/middleware/lrw"

	"github.com/promoboxx/go-glitch/glitch"
	"github.com/promoboxx/go-service/alice/middleware/contextkey"
	"github.com/sirupsen/logrus"
)

const (
	logFieldRequestID   = "request_id"
	logFieldUserID      = "user_id"
	logFieldMethod      = "method"
	logFieldStatusCode  = "status_code"
	logFieldPath        = "path"
	logFieldError       = "error"
	logFieldTraceID     = "dd.trace_id"
	logFieldQueryParams = "query_params"
)

// Logger injects a logger into the context
type Logger interface {
	Log(h http.Handler) http.Handler
}

type logger struct {
	entry       *logrus.Entry
	logRequests bool
}

// NewLogrusLogger allows you to setup a base entry to use for logging
func NewLogrusLogger(baseEntry *logrus.Entry, logRequests bool) Logger {
	return &logger{entry: baseEntry, logRequests: logRequests}
}

// Log is the middleware for injecting a logger into the context that has
// request specific information
func (l *logger) Log(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// base fields that get added to each log entry
		fields := logrus.Fields{
			logFieldUserID:      GetInsecureUserIDFromContext(r.Context()),
			logFieldMethod:      r.Method,
			logFieldPath:        r.URL.Path,
			logFieldQueryParams: r.URL.Query(),
			// this header comes from the data dog span information that gets injected
			// every request from the headers
			logFieldTraceID: r.Header.Get("X-Datadog-Trace-ID"),
		}

		entry := l.entry.WithFields(fields)

		// add logger to the context
		ctx := context.WithValue(r.Context(), contextkey.ContextKeyLogger, entry)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)

		responseFields := fields

		if loggingResponseWriter, ok := w.(*lrw.LoggingResponseWriter); ok {
			responseFields[logFieldStatusCode] = loggingResponseWriter.StatusCode

			if loggingResponseWriter.InnerError != nil {
				fields[logFieldError] = loggingResponseWriter.InnerError

				if dataErr, ok := loggingResponseWriter.InnerError.(glitch.DataError); ok {
					for k, v := range dataErr.GetFields() {
						fields[k] = v
					}
				}
			}

			for fieldName, message := range loggingResponseWriter.ExtraFields {
				fields[fieldName] = message
			}
		}

		responseEntry := l.entry.WithFields(responseFields)

		if l.logRequests {
			responseEntry.Printf("Finished request")
		}
	})
}
