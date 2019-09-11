package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_CreateLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	var isFailed bool
	lb := getMockLoadbalancer().Properties
	labelSlices := [][]string{nil, lb.Labels}
	uri := path.Join(apiLoadBalancerBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		if isFailed {
			w.WriteHeader(400)
		} else {
			fmt.Fprint(w, prepareLoadBalancerHTTPCreateResponse())
		}
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	for _, testSuccessFail := range commonSuccessFailTestCases {
		isFailed = testSuccessFail.isFailed
		for _, testLabel := range labelSlices {
			lbRequest := LoadBalancerCreateRequest{
				Name:                lb.Name,
				Algorithm:           lb.Algorithm,
				LocationUUID:        lb.LocationUUID,
				ListenIPv6UUID:      lb.ListenIPv6UUID,
				ListenIPv4UUID:      lb.ListenIPv4UUID,
				RedirectHTTPToHTTPS: lb.RedirectHTTPToHTTPS,
				ForwardingRules:     lb.ForwardingRules,
				BackendServers:      lb.BackendServers,
				Labels:              testLabel,
			}
			response, err := client.CreateLoadBalancer(lbRequest)
			if testSuccessFail.isFailed {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err, "CreateLoadBalancer returned error: %v", err)
				assert.Equal(t, fmt.Sprintf("%s", prepareLoadBalancerObjectCreateResponse()), fmt.Sprintf("%s", response))
			}
		}
	}
}

func TestClient_GetLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	expectedObject := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPGetResponse())
	})
	for _, test := range uuidCommonTestCases {
		loadbalancer, err := client.GetLoadBalancer(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLoadBalancer returned error: %v", err)
			assert.Equal(t, fmt.Sprintf("%v", expectedObject.Properties), fmt.Sprintf("%v", loadbalancer.Properties))
		}
	}
}

func TestClient_GetLoadBalancerList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase)
	expectedObjects := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareLoadBalancerHTTPListResponse())
	})
	loadbalancers, err := client.GetLoadBalancerList()
	assert.Nil(t, err, "GetLoadBalancerList returned error: %v", err)
	assert.Equal(t, 1, len(loadbalancers))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", loadbalancers))
}

func TestClient_UpdateLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.UpdateLoadBalancer(test.testUUID, LoadBalancerUpdateRequest{
			Name:                "test",
			ListenIPv6UUID:      dummyUUID,
			ListenIPv4UUID:      dummyUUID,
			RedirectHTTPToHTTPS: false,
			Status:              "inactive",
		})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "UpdateLoadBalancer returned an error %v", err)
		}
	}
}

func TestClient_DeleteLoadBalancer(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.DeleteLoadBalancer(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteLoadBalancer returned an error %v", err)
		}
	}
}

func TestClient_GetLoadBalancerEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLoadBalancerBase, dummyUUID, "events")
	expectedObjects := getMockEvent()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetLoadBalancerEventList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetLoadBalancerEventList returned error: %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", response))
		}
	}
}

func getMockLoadbalancer() LoadBalancer {
	labels := make([]string, 0)
	labels = append(labels, "nice")
	lb := LoadBalancer{
		Properties: LoadBalancerProperties{
			ObjectUUID:          dummyUUID,
			Name:                "go-client-lb",
			Algorithm:           "leastconn",
			LocationUUID:        "45ed677b-3702-4b36-be2a-a2eab9827950",
			ListenIPv6UUID:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			ListenIPv4UUID:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			RedirectHTTPToHTTPS: false,
			ForwardingRules: []ForwardingRule{
				{
					LetsencryptSSL: "",
					ListenPort:     8080,
					Mode:           "http",
					TargetPort:     8000,
				},
			},
			BackendServers: []BackendServer{
				{
					Weight: 100,
					Host:   "185.201.147.176",
				},
			},
			Labels: labels,
		},
	}
	return lb
}

func prepareLoadBalancerHTTPGetResponse() string {
	lb := getMockLoadbalancer()
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancer": %s}`, string(res))
}

func prepareLoadBalancerHTTPListResponse() string {
	lb := getMockLoadbalancer()
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancers": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareLoadBalancerHTTPCreateResponse() string {
	return fmt.Sprintf(`{"request_uuid": "%s","object_uuid": "%s"}`, dummyRequestUUID, dummyUUID)
}

func prepareLoadBalancerObjectCreateResponse() LoadBalancerCreateResponse {
	return LoadBalancerCreateResponse{
		RequestUUID: dummyRequestUUID,
		ObjectUUID:  dummyUUID,
	}
}
