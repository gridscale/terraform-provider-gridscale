package gsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"gotest.tools/assert"
)

const loadBalancerID = "690de890-13c0-4e76-8a01-e10ba8786e53"
const requestUUID = "x123xx1x-123x-1x12-123x-123xxx123x1x"

func setupTestClient() (*Client, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	config := Config{
		APIUrl:     server.URL,
		UserUUID:   "uuid",
		APIToken:   "token",
		HTTPClient: http.DefaultClient,
	}
	return NewClient(&config), mux
}

func TestCreateLoadBalancer(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiLoadBalancerBase)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		fmt.Fprint(w, prepareHTTPCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, requestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	lb := getMockLoadbalancer().Properties
	lbRequest := LoadBalancerCreateRequest{
		Name:                lb.Name,
		Algorithm:           lb.Algorithm,
		LocationUuid:        lb.LocationUuid,
		ListenIPv6Uuid:      lb.ListenIPv6Uuid,
		ListenIPv4Uuid:      lb.ListenIPv4Uuid,
		RedirectHTTPToHTTPS: lb.RedirectHTTPToHTTPS,
		ForwardingRules:     lb.ForwardingRules,
		BackendServers:      lb.BackendServers,
		Labels:              lb.Labels,
	}
	response, err := client.CreateLoadBalancer(lbRequest)
	if err != nil {
		t.Errorf("CreateLoadBalancer returned error: %v", err)
	}
	assert.Equal(t, fmt.Sprintf("&%s", prepareObjectCreateResponse()), fmt.Sprintf("%s", response))
}
func TestGetLoadBalancer(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiLoadBalancerBase, loadBalancerID)
	expectedObject := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareHTTPGetResponse())
	})
	loadbalancer, err := client.GetLoadBalancer(loadBalancerID)
	if err != nil {
		t.Errorf("GetLoadBalancer returned error: %v", err)
	}
	assert.Equal(t, fmt.Sprintf("%v", expectedObject.Properties), fmt.Sprintf("%v", loadbalancer.Properties))
}
func TestGetLoadBalancerList(t *testing.T) {
	client, mux := setupTestClient()
	uri := path.Join(apiLoadBalancerBase)
	expectedObjects := getMockLoadbalancer()
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodGet)
		fmt.Fprint(w, prepareHTTPListResponse())
	})
	loadbalancers, err := client.GetLoadBalancerList()
	if err != nil {
		t.Errorf("GetLoadBalancerList returned error: %v", err)
	}
	assert.Equal(t, 1, len(loadbalancers))
	assert.Equal(t, fmt.Sprintf("[%v]", expectedObjects), fmt.Sprintf("%v", loadbalancers))
}

func getMockLoadbalancer() LoadBalancer {
	labels := make([]interface{}, 0)
	labels = append(labels, "nice")
	lb := LoadBalancer{
		Properties: LoadBalancerProperties{
			ObjectUuid:          loadBalancerID,
			Name:                "go-client-lb",
			Algorithm:           "leastconn",
			LocationUuid:        "45ed677b-3702-4b36-be2a-a2eab9827950",
			ListenIPv6Uuid:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			ListenIPv4Uuid:      "880b7f98-3702-4b36-be2a-a2eab9827950",
			RedirectHTTPToHTTPS: false,
			ForwardingRules: []ForwardingRule{
				{
					LetsencryptSSL: nil,
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

func prepareHTTPGetResponse() string {
	lb := getMockLoadbalancer()
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancer": %s}`, string(res))
}

func prepareHTTPListResponse() string {
	lb := getMockLoadbalancer()
	res, _ := json.Marshal(lb.Properties)
	return fmt.Sprintf(`{"loadbalancers": {"%s": %s}}`, loadBalancerID, string(res))
}

func prepareHTTPCreateResponse() string {
	return fmt.Sprintf(`{"request_uuid": "%s","object_uuid": "%s"}`, requestUUID, loadBalancerID)
}

func prepareObjectCreateResponse() LoadBalancerCreateResponse {
	return LoadBalancerCreateResponse{
		RequestUuid: requestUUID,
		ObjectUuid:  loadBalancerID}
}
