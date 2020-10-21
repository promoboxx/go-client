package lrw

import (
	"fmt"
	"net/http"
)

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode  int
	InnerError  error
	ExtraFields map[string]string
}

type InvalidFieldError struct {
	InvalidField string
	Message      string
}

func (l *LoggingResponseWriter) WriteHeader(code int) {
	l.StatusCode = code
	l.ResponseWriter.WriteHeader(code)
}

func NewLoggingResponseWriter(rw http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{ResponseWriter: rw, StatusCode: http.StatusOK, InnerError: nil, ExtraFields: map[string]string{}}
}

// AddLogField will attempt to add the field to the logs that will emitted for each request, it will fail if it attempts to override another field
func (l *LoggingResponseWriter) AddLogField(name string, value string) error {
	if _, ok := l.ExtraFields[name]; !ok {
		l.ExtraFields[name] = value
	} else {
		return newInvalidFieldError(name)
	}

	return nil
}

// ForceAddLogField will always add the log field and will overwrite anything that may currently exist as a field already
func (l *LoggingResponseWriter) ForceAddLogField(name string, value string) {
	l.ExtraFields[name] = value
}

func (ife *InvalidFieldError) Error() string {
	return ife.Message
}

func newInvalidFieldError(name string) *InvalidFieldError {
	return &InvalidFieldError{InvalidField: name, Message: fmt.Sprintf("field %s overrides an already existing field", name)}
}

// AddLogField will attempt to add the log field to the response writer, if it is not a
// LoggingResponseWriter then it will turn the ResponseWriter into one
func AddLogField(rw http.ResponseWriter, name string, value string) error {
	if lrw, ok := rw.(*LoggingResponseWriter); ok {
		return lrw.AddLogField(name, value)
	}

	lrw := NewLoggingResponseWriter(rw)
	err := lrw.AddLogField(name, value)
	if err != nil {
		return err
	}
	rw = lrw

	return nil
}

// ForceAddLogField will always to add the log field to the response writer and ovewrite any field that may already exist
// if it is not a LoggingResponseWriter then it will turn the ResponseWriter into one
func ForceAddLogField(rw http.ResponseWriter, name string, value string) {
	if lrw, ok := rw.(*LoggingResponseWriter); ok {
		lrw.ForceAddLogField(name, value)
		return
	}

	lrw := NewLoggingResponseWriter(rw)
	lrw.ForceAddLogField(name, value)
	rw = lrw
}
