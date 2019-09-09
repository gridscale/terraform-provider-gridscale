package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetLocationList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiLocationBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareLocationListHTTPGet())
	})
	res, err := client.GetLocationList()
	assert.Nil(t, err, "GetLocationList returned an aerror %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockLocation()), fmt.Sprintf("%v", res))
}

func TestClient_GetLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareLocationHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetLocation(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLocation returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockLocation()), fmt.Sprintf("%v", res))
		}
	}
}

func getMockLocation() Location {
	mock := Location{Properties: LocationProperties{
		Iata:       "fra",
		Status:     "active",
		Labels:     nil,
		Name:       "de/fra",
		ObjectUUID: dummyUUID,
		Country:    "de",
	}}
	return mock
}

func prepareLocationHTTPGet() string {
	location := getMockLocation()
	res, _ := json.Marshal(location)
	return string(res)
}

func prepareLocationListHTTPGet() string {
	location := getMockLocation()
	res, _ := json.Marshal(location.Properties)
	return fmt.Sprintf(`{"locations": {"%s": %s}}`, dummyUUID, string(res))
}
