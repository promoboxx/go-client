package lrw

import "net/http"

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	InnerError error
}

func (l *LoggingResponseWriter) WriteHeader(code int) {
	l.StatusCode = code
	l.ResponseWriter.WriteHeader(code)
}
