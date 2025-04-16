package glitch

import "fmt"

// UnknownCode is used as the code when a real code could not be found
const UnknownCode = "UNKNOWN"

// DataError is a class of error that provides access to an error code as well as the originating error
type DataError interface {
	// Error satisfies the error interface
	Error() string
	// Inner returns the originating error
	Inner() error
	// Code returns a specific error code, this is meant to be a machine friendly string
	Code() string
	// Msg returns the originating message
	Msg() string
	// Wrap will set err as the cause of the error and returns itself
	Wrap(err DataError) DataError
	// GetCause will return the cause of this error
	GetCause() DataError
	// IsTransient will return if the error is considered transient
	IsTransient() bool
	// GetFields returns the extra fields for giving better error descriptions
	GetFields() map[string]interface{}
	// GetField returns the data of a field from the fields member of the DataError if one exists
	GetField(name string) interface{}
	// AddFields will add fields to the given DataErrors fields
	AddFields(map[string]interface{})
	// AddField will add the key and value pair to the data errors fields
	AddField(key string, value interface{})
	// String implements the string interface for %s in fmt calls
	String() string
}

type dataError struct {
	inner       error
	code        string
	msg         string
	cause       DataError
	isTransient bool
	fields      map[string]interface{}
}

func (d *dataError) Error() string {
	return fmt.Sprintf("Code: [%s] Message: [%s] Inner error: [%s]", d.code, d.msg, d.inner)
}

func (d *dataError) Inner() error {
	return d.inner
}

func (d *dataError) Code() string {
	return d.code
}

func (d *dataError) Msg() string {
	return d.msg
}

func (d *dataError) Wrap(err DataError) DataError {
	d.cause = err
	return d
}

func (d *dataError) GetCause() DataError {
	if d.cause == nil {
		return d
	}

	return d.cause
}

func (d *dataError) IsTransient() bool {
	return d.isTransient
}

func (d *dataError) GetFields() map[string]interface{} {
	return d.fields
}

func (d *dataError) GetField(name string) interface{} {
	return d.fields[name]
}

// AddFields adds the given fields to the data error
// NOTE: Any field given that matches one in the data DataErrors
// fields already will overwrite it
func (d *dataError) AddFields(fields map[string]interface{}) {
	for k, v := range fields {
		d.AddField(k, v)
	}
}

func (d *dataError) AddField(key string, value interface{}) {
	d.fields[key] = value
}

func (d *dataError) String() string {
	return d.Error()
}

// FromHTTPProblem will create a DataError from an HTTPProblem
func FromHTTPProblem(inner error, msg string) DataError {
	switch ie := inner.(type) {
	case HTTPProblem:
		return &dataError{inner: inner, code: ie.Code, msg: msg, isTransient: ie.IsTransient}
	case HTTPProblemMetadata:
		return &dataError{inner: inner, code: ie.Code, msg: msg, isTransient: ie.IsTransient, fields: map[string]interface{}{"metadata": ie.Metadata}}
	}

	return &dataError{inner: inner, code: UnknownCode, msg: msg}
}

// NewDataError will create a DataError from the information provided
func NewDataError(inner error, code string, msg string) DataError {
	return &dataError{inner: inner, code: code, msg: msg, fields: map[string]interface{}{}}
}

// NewTransientDataError will create a DataError from the information provided and set the isTransient flag to true
func NewTransientDataError(inner error, code string, msg string) DataError {
	return &dataError{inner: inner, code: code, msg: msg, isTransient: true, fields: map[string]interface{}{}}
}
