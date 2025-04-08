package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/promoboxx/go-glitch/glitch"
	"github.com/promoboxx/go-service/alice/middleware/lrw"
)

// ReturnProblem will return a json http problem response
func ReturnProblem(w http.ResponseWriter, detail, code string, status int, innerErr error) (int, []byte) {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}

	if dataErr, ok := innerErr.(glitch.DataError); ok {
		prob.IsTransient = dataErr.IsTransient()
	}

	if loggingResponseWriter, ok := w.(*lrw.LoggingResponseWriter); ok {
		loggingResponseWriter.InnerError = innerErr
	}

	by, _ := json.Marshal(prob)
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}

	return status, by
}

// WriteProblem will write a json http problem response
func WriteProblem(w http.ResponseWriter, detail, code string, status int, innerErr error) error {
	prob := glitch.HTTPProblem{
		Title:  http.StatusText(status),
		Detail: detail,
		Code:   code,
		Status: status,
	}

	if dataErr, ok := innerErr.(glitch.DataError); ok {
		prob.IsTransient = dataErr.IsTransient()
	}

	by, err := json.Marshal(prob)
	if err != nil {
		return err
	}

	if lrw, ok := w.(*lrw.LoggingResponseWriter); ok {
		lrw.InnerError = innerErr
		lrw.AddLogField("error_code", code)
		lrw.AddLogField("error_detail", detail)
	}

	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// WriteProblemWithMetadata writes a normal http problem but will also add metadata to the response
func WriteProblemWithMetadata(w http.ResponseWriter, detail, code string, status int, innerErr error, metadata interface{}) error {
	prob := glitch.HTTPProblemMetadata{
		HTTPProblem: glitch.HTTPProblem{
			Title:  http.StatusText(status),
			Detail: detail,
			Code:   code,
			Status: status,
		},
		Metadata: metadata,
	}

	if dataErr, ok := innerErr.(glitch.DataError); ok {
		prob.IsTransient = dataErr.IsTransient()
	}

	by, err := json.Marshal(prob)
	if err != nil {
		return err
	}

	if lrw, ok := w.(*lrw.LoggingResponseWriter); ok {
		lrw.InnerError = innerErr
		lrw.AddLogField("error_code", code)
		lrw.AddLogField("error_detail", detail)
	}

	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// Writes the given error as a json http problem response to the `http.ResponseWriter`, and
// returns the raw error
func WriteDataError(w http.ResponseWriter, err glitch.DataError, status int) error {
	return WriteProblem(w, err.Msg(), err.Code(), status, err.Inner())
}

// WriteJSONResponse will write a json response to the htt.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	var by []byte
	var err error
	if data != nil {
		by, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	if w != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		_, err = w.Write(by)
	}
	return err
}

// Int32PointerFromQueryParam returns a nullable int32 from a query param key
func Int32PointerFromQueryParam(r *http.Request, paramName string) (*int32, error) {
	strValue := r.URL.Query().Get(paramName)
	var intPointer *int32
	if len(strValue) > 0 {
		i, err := strconv.ParseInt(strValue, 10, 32)
		if err != nil {
			return intPointer, err
		}
		i32 := int32(i)
		intPointer = &i32
	}
	return intPointer, nil
}

func Int64ArrayFromQueryParam(r *http.Request, paramName string) ([]int64, error) {
	var ret []int64
	str := r.URL.Query().Get(paramName)
	if len(str) == 0 {
		return ret, nil
	}
	parts := strings.Split(str, ",")
	for _, v := range parts {
		i, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("(%s) from input (%s) was not an integer: %s", v, str, err)
		}
		ret = append(ret, i)
	}
	return ret, nil
}

func Int32ArrayFromQueryParam(r *http.Request, paramName string) ([]int32, error) {
	var result []int32
	str := r.URL.Query().Get(paramName)
	if len(str) == 0 {
		return result, nil
	}
	parts := strings.Split(str, ",")

	for _, v := range parts {
		i64, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("%s is not an integer", v)
		}
		result = append(result, int32(i64))
	}
	return result, nil
}

func TimestampFromQueryParam(r *http.Request, paramName string) (*time.Time, error) {
	str := r.URL.Query().Get(paramName)
	if len(str) == 0 {
		return nil, nil
	}
	ret, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return nil, fmt.Errorf("(%s) was not a valid timestamp: %s", str, err)
	}
	return &ret, nil
}
