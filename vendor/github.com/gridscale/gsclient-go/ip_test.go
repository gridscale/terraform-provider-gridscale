package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetIPList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPListHTTPGet())
	})
	res, err := client.GetIPList()
	assert.Nil(t, err, "GetIPList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
}

func TestClient_GetIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetIP(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetIP returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockIP()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateIP(t *testing.T) {
	server, client, mux := setupTestClient()
	var isFailed bool
	defer server.Close()
	uri := apiIPBase
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareIPCreateResponse())
		}
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		response, err := client.CreateIP(IPCreateRequest{
			Name:         "test",
			Family:       1,
			LocationUUID: dummyUUID,
			Failover:     false,
			ReverseDNS:   "8.8.8.8",
		})
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "CreateIP returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockIPCreateResponse()), fmt.Sprintf("%s", response))
		}
	}
}

func TestClient_UpdateIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPatch, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.UpdateIP(test.testUUID, IPUpdateRequest{
			Name:       "test",
			Failover:   false,
			ReverseDNS: "8.8.4.4",
		})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "UpdateIP returned an error %v", err)
		}
	}
}

func TestClient_DeleteIP(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.DeleteIP(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteIP returned an error %v", err)
		}
	}
}

func TestClient_GetIPEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiIPBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetIPEventList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetIPEventList returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", res))
		}
	}

}

func TestClient_GetIPVersion(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiIPBase, dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprintf(writer, prepareIPHTTPGet())
		}
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res := client.GetIPVersion(dummyUUID)
		if test.isFailed {
			assert.Equal(t, 0, res)
		} else {
			assert.Equal(t, 1, res)
		}
	}
}

func TestClient_GetIPsByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareIPListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetIPsByLocation(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetIPsByLocation returned an error %v", err)
			assert.Equal(t, 1, len(res))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_GetDeletedIPs(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiDeletedBase, "ips")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareDeletedIPListHTTPGet())
	})
	res, err := client.GetDeletedIPs()
	assert.Nil(t, err, "GetDeletedIPs returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockIP()), fmt.Sprintf("%v", res))
}

func getMockIP() IP {
	mock := IP{Properties: IPProperties{
		Name:            "test",
		LocationCountry: "Germany",
		LocationUUID:    dummyUUID,
		ObjectUUID:      dummyUUID,
		ReverseDNS:      "8.8.8.8",
		Family:          1,
		Status:          "active",
		CreateTime:      dummyTime,
		Failover:        false,
		ChangeTime:      dummyTime,
		LocationIata:    "",
		LocationName:    "Cologne",
		Prefix:          "",
		IP:              "192.168.0.1",
		DeleteBlock:     "",
		UsagesInMinutes: 10,
		CurrentPrice:    0.9,
		Labels:          []string{"label"},
		Relations: IPRelations{
			Loadbalancers: []IPLoadbalancer{
				{
					CreateTime:       dummyTime,
					LoadbalancerName: "test",
					LoadbalancerUUID: dummyUUID,
				},
			},
		},
	}}
	return mock
}

func prepareIPListHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"ips": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareIPHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip)
	return string(res)
}

func getMockIPCreateResponse() IPCreateResponse {
	mock := IPCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
		Prefix:      "ip",
		IP:          "192.168.0.1",
	}
	return mock
}

func prepareIPCreateResponse() string {
	res, _ := json.Marshal(getMockIPCreateResponse())
	return string(res)
}

func prepareDeletedIPListHTTPGet() string {
	ip := getMockIP()
	res, _ := json.Marshal(ip.Properties)
	return fmt.Sprintf(`{"deleted_ips": {"%s": %s}}`, dummyUUID, string(res))
}
