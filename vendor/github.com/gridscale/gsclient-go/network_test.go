package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetNetworkList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet(true))
	})
	res, err := client.GetNetworkList()
	assert.Nil(t, err, "GetNetworkList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true)), fmt.Sprintf("%v", res))
}

func TestClient_GetNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetwork(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetwork returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockNetwork(true)), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	var isFailed bool
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareNetworkCreateResponse())
		}
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreateNetwork(NetworkCreateRequest{
			Name:         "test",
			Labels:       []string{"label"},
			LocationUUID: dummyUUID,
			L2Security:   false,
		})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateNetwork returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockNetworkCreateResponse()), fmt.Sprintf("%s", response))
		}
	}
}

func TestClient_UpdateNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.UpdateNetwork(test.testUUID, NetworkUpdateRequest{
			Name:       "test",
			L2Security: false,
		})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "UpdateNetwork returned an error %v", err)
		}
	}
}

func TestClient_DeleteNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.DeleteNetwork(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteNetwork returned an error %v", err)
		}
	}
}

func TestClient_GetNetworkEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiNetworkBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetworkEventList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetworkEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetNetworkPublic(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	var isFailed bool
	var isPublicNet bool
	pubNetCases := []bool{true, false}
	uri := apiNetworkBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareNetworkListHTTPGet(isPublicNet))
		}
	})
	for _, successFailTest := range commonSuccessFailTestCases {
		isFailed = successFailTest.isFailed
		for _, publicNetTest := range pubNetCases {
			isPublicNet = publicNetTest
			res, err := client.GetNetworkPublic()
			if isFailed || !publicNetTest {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetNetworkPublic returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockNetwork(publicNetTest)), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_GetNetworksByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareNetworkListHTTPGet(true))
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetNetworksByLocation(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetNetworksByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true)), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetDeletedNetworks(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiDeletedBase, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareDeletedNetworkListHTTPGet())
	})
	res, err := client.GetDeletedNetworks()
	assert.Nil(t, err, "GetDeletedNetworks returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockNetwork(true)), fmt.Sprintf("%v", res))
}

func getMockNetwork(isPublic bool) Network {
	mock := Network{Properties: NetworkProperties{
		LocationCountry: "Germany",
		LocationUUID:    "",
		PublicNet:       isPublic,
		ObjectUUID:      dummyUUID,
		NetworkType:     "",
		Name:            "test",
		Status:          "active",
		CreateTime:      dummyTime,
		L2Security:      false,
		ChangeTime:      dummyTime,
		LocationName:    "Cologne",
		DeleteBlock:     false,
		Labels:          nil,
		Relations: NetworkRelations{
			Vlans: []NetworkVlan{
				{
					Vlan:       1,
					TenantName: "test",
					TenantUUID: dummyUUID,
				},
			},
		},
	}}
	return mock
}

func prepareNetworkListHTTPGet(isPublic bool) string {
	network := getMockNetwork(isPublic)
	res, _ := json.Marshal(network.Properties)
	return fmt.Sprintf(`{"networks": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareNetworkHTTPGet() string {
	network := getMockNetwork(true)
	res, _ := json.Marshal(network)
	return string(res)
}

func getMockNetworkCreateResponse() NetworkCreateResponse {
	mock := NetworkCreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareNetworkCreateResponse() string {
	createResponse := getMockNetworkCreateResponse()
	res, _ := json.Marshal(createResponse)
	return string(res)
}

func prepareDeletedNetworkListHTTPGet() string {
	network := getMockNetwork(true)
	res, _ := json.Marshal(network.Properties)
	return fmt.Sprintf(`{"deleted_networks": {"%s": %s}}`, dummyUUID, string(res))
}
