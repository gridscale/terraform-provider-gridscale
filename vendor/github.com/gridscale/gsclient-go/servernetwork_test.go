package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetServerNetworkList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetServerNetworkList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetServerNetworkList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockServerNetwork()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareServerNetworkHTTPGet())
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testNetworkID := range uuidCommonTestCases {
			res, err := client.GetServerNetwork(testServerID.testUUID, testNetworkID.testUUID)
			if testServerID.isFailed || testNetworkID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "GetServerNetwork returned an error %v", err)
				assert.Equal(t, fmt.Sprintf("%v", getMockServerNetwork()), fmt.Sprintf("%v", res))
			}
		}
	}
}

func TestClient_CreateServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testNetworkID := range uuidCommonTestCases {
			err := client.CreateServerNetwork(testServerID.testUUID, ServerNetworkRelationCreateRequest{
				ObjectUUID:           testNetworkID.testUUID,
				Ordering:             1,
				BootDevice:           false,
				L3security:           nil,
				FirewallTemplateUUID: dummyUUID,
			})
			if testServerID.isFailed || testNetworkID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateServerNetwork returned an error %v", err)
			}
		}
	}
}

func TestClient_UpdateServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testNetworkID := range uuidCommonTestCases {
			err := client.UpdateServerNetwork(testServerID.testUUID, testNetworkID.testUUID, ServerNetworkRelationUpdateRequest{
				Ordering:             0,
				BootDevice:           true,
				FirewallTemplateUUID: dummyUUID,
			})
			if testServerID.isFailed || testNetworkID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "UpdateServerNetwork returned an error %v", err)
			}
		}
	}
}

func TestClient_DeleteServerNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, testServerID := range uuidCommonTestCases {
		for _, testNetworkID := range uuidCommonTestCases {
			err := client.DeleteServerNetwork(testServerID.testUUID, testNetworkID.testUUID)
			if testServerID.isFailed || testNetworkID.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "DeleteServerNetwork returned an error %v", err)
			}
		}
	}
}

func TestClient_LinkNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.LinkNetwork(dummyUUID, dummyUUID, dummyUUID, true, 0, nil, FirewallRules{})
	assert.Nil(t, err, "LinkNetwork returned an error %v", err)
}

func TestClient_UnlinkNetwork(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiServerBase, dummyUUID, "networks", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	err := client.UnlinkNetwork(dummyUUID, dummyUUID)
	assert.Nil(t, err, "UnlinkNetwork returned an error %v", err)
}

func getMockServerNetwork() ServerNetworkRelationProperties {
	mock := ServerNetworkRelationProperties{
		L2security:           true,
		ServerUUID:           dummyUUID,
		CreateTime:           dummyTime,
		PublicNet:            false,
		FirewallTemplateUUID: dummyUUID,
		ObjectName:           "test",
		Mac:                  "",
		BootDevice:           true,
		PartnerUUID:          dummyUUID,
		Ordering:             0,
		Firewall:             "",
		NetworkType:          "",
		NetworkUUID:          dummyUUID,
		ObjectUUID:           dummyUUID,
		L3security:           nil,
	}
	return mock
}

func prepareServerNetworkListHTTPGet() string {
	net := getMockServerNetwork()
	res, _ := json.Marshal(net)
	return fmt.Sprintf(`{"network_relations": [%s]}`, string(res))
}

func prepareServerNetworkHTTPGet() string {
	net := getMockServerNetwork()
	res, _ := json.Marshal(net)
	return fmt.Sprintf(`{"network_relation": %s}`, string(res))
}
