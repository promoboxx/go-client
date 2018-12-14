package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/promoboxx/go-glitch/glitch"
)

// ObjectToJSONReader will v to a io.Reader of the JSON representation of v
func ObjectToJSONReader(v interface{}) (io.Reader, glitch.DataError) {
	if by, ok := v.([]byte); ok {
		return bytes.NewBuffer(by), nil
	}
	by, err := json.Marshal(v)
	if err != nil {
		return nil, glitch.NewDataError(err, ErrorMarshallingObject, "Error marshalling object to json")
	}
	return bytes.NewBuffer(by), nil
}

// BuildRoute is a simple helper function to build the route for service to service communication
func BuildRoute(slug string, appendServiceNameToRoute bool, pathPrefix string) string {
	if appendServiceNameToRoute {
		return fmt.Sprintf("%s/%s", pathPrefix, slug)
	}

	return slug
}
