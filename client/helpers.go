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

// PrefixRoute is a simple helper function to build the route for service to service communicationfunc PrefixRoute(pathPrefix string, appendServiceNameToRoute bool, route string) string {
func PrefixRoute(serviceName string, pathPrefix string, appendServiceNameToRoute bool, route string) string {
	routePrefix := fmt.Sprintf("%s", pathPrefix)
	if appendServiceNameToRoute {
		routePrefix = fmt.Sprintf("%s/%s", pathPrefix, serviceName)
	}
	// TODO: make this more robust by normalizing slashes
	return fmt.Sprintf("%s/%s", routePrefix, route)
}
