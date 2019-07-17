package gsclient

import (
	"path"
)

//LoadBalancers is the JSON struct of a list of loadbalancers
type LoadBalancers struct {
	List map[string]LoadBalancerProperties `json:"loadbalancers"`
}

//LoadBalancer is the JSON struct of a loadbalancer
type LoadBalancer struct {
	Properties LoadBalancerProperties `json:"loadbalancer"`
}

//LoadBalancerProperties is the properties of a loadbalancer
type LoadBalancerProperties struct {
	ObjectUuid          string                    `json:"object_uuid"`
	Name                string                    `json:"name"`
	ListenIPv6Uuid      string                    `json:"listen_ipv6_uuid"`
	ListenIPv4Uuid      string                    `json:"listen_ipv4_uuid"`
	Algorithm           string                    `json:"algorithm"`
	ForwardingRules     []ForwardingRule          `json:"forwarding_rules"`
	BackendServers      []BackendServer           `json:"backend_servers"`
	Labels              []string                  `json:"labels"`
	LocationUuid        string                    `json:"location_uuid"`
	RedirectHTTPToHTTPS bool                      `json:"redirect_http_to_https"`
	CreateTime          string                    `json:"create_time"`
	ListenPort          map[string]map[string]int `json:"listen_ports"`
	ChangeTime          string                    `json:"change_time"`
	Status              string                    `json:"status"`
}

//BackendServer is the JSON struct of backend server
type BackendServer struct {
	Weight int    `json:"weight"`
	Host   string `json:"host"`
}

//ForwardingRule is the JSON struct of forwarding rule
type ForwardingRule struct {
	LetsencryptSSL *string `json:"letsencrypt_ssl"`
	ListenPort     int     `json:"listen_port"`
	Mode           string  `json:"mode"`
	TargetPort     int     `json:"target_port"`
}

//LoadBalancerCreateRequest is the JSON struct for creating a loadbalancer request
type LoadBalancerCreateRequest struct {
	Name                string           `json:"name"`
	ListenIPv6Uuid      string           `json:"listen_ipv6_uuid"`
	ListenIPv4Uuid      string           `json:"listen_ipv4_uuid"`
	Algorithm           string           `json:"algorithm"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	Labels              []interface{}    `json:"labels"`
	LocationUuid        string           `json:"location_uuid"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Status              string           `json:"status,omitempty"`
}

//LoadBalancerUpdateRequest is the JSON struct for updating a loadbalancer request
type LoadBalancerUpdateRequest struct {
	Name                string           `json:"name"`
	ListenIPv6Uuid      string           `json:"listen_ipv6_uuid"`
	ListenIPv4Uuid      string           `json:"listen_ipv4_uuid"`
	Algorithm           string           `json:"algorithm"`
	ForwardingRules     []ForwardingRule `json:"forwarding_rules"`
	BackendServers      []BackendServer  `json:"backend_servers"`
	Labels              []string         `json:"labels"`
	LocationUuid        string           `json:"location_uuid"`
	RedirectHTTPToHTTPS bool             `json:"redirect_http_to_https"`
	Status              string           `json:"status,omitempty"`
}

//LoadBalancerCreateResponse is the JSON struct for a loadbalancer response
type LoadBalancerCreateResponse struct {
	RequestUuid string `json:"request_uuid"`
	ObjectUuid  string `json:"object_uuid"`
}

//LoadBalancerEvents is the JSON struct for a loadbalancer's events
type LoadBalancerEvents struct {
	Events []LoadBalancerEventProperties `json:"events"`
}

//LoadBalancerEventProperties is the properties of a loadbalancer's event
type LoadBalancerEventProperties struct {
	ObjectUuid    string `json:"object_uuid"`
	ObjectType    string `json:"object_type"`
	RequestUuid   string `json:"request_uuid"`
	RequestType   string `json:"request_type"`
	Activity      string `json:"activity"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUuid      string `json:"user_uuid"`
}

//GetLoadBalancerList returns a list of load balancers
func (c *Client) GetLoadBalancerList() ([]LoadBalancer, error) {
	r := Request{
		uri:    apiLoadBalancerBase,
		method: "GET",
	}

	response := new(LoadBalancers)
	err := r.execute(*c, &response)

	list := []LoadBalancer{}
	for _, properties := range response.List {
		list = append(list, LoadBalancer{Properties: properties})
	}

	return list, err
}

//GetLoadBalancer returns a load balancer of a given uuid
func (c *Client) GetLoadBalancer(id string) (*LoadBalancer, error) {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: "GET",
	}

	response := new(LoadBalancer)
	err := r.execute(*c, &response)

	return response, err
}

//CreateLoadBalancer creates a new load balancer
func (c *Client) CreateLoadBalancer(body LoadBalancerCreateRequest) (*LoadBalancerCreateResponse, error) {
	r := Request{
		uri:    apiLoadBalancerBase,
		method: "POST",
		body:   body,
	}

	response := new(LoadBalancerCreateResponse)
	err := r.execute(*c, &response)
	if err != nil {
		return nil, err
	}

	err = c.WaitForRequestCompletion(response.RequestUuid)

	return response, err
}

//UpdateLoadBalancer update configuraton of a load balancer
func (c *Client) UpdateLoadBalancer(id string, body LoadBalancerUpdateRequest) error {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}

//GetLoadBalancerEventList retrives events of a given uuid
func (c *Client) GetLoadBalancerEventList(id string) (*LoadBalancerEvents, error) {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id, "events"),
		method: "GET",
	}
	response := new(LoadBalancerEvents)
	err := r.execute(*c, &response)
	return response, err
}

//DeleteLoadBalancer deletes a load balancer
func (c *Client) DeleteLoadBalancer(id string) error {
	r := Request{
		uri:    path.Join(apiLoadBalancerBase, id),
		method: "DELETE",
	}

	return r.execute(*c, nil)
}
