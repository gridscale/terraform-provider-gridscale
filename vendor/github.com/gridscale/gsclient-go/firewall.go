package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//FirewallList is JSON structure of a list of firewalls
type FirewallList struct {
	List map[string]FirewallProperties `json:"firewalls"`
}

//Firewall is JSON structure of a single firewall
type Firewall struct {
	Properties FirewallProperties `json:"firewall"`
}

//FirewallProperties is JSON struct of a firewall's properties
type FirewallProperties struct {
	Status       string           `json:"status"`
	Labels       []string         `json:"labels"`
	ObjectUUID   string           `json:"object_uuid"`
	ChangeTime   string           `json:"change_time"`
	Rules        FirewallRules    `json:"rules"`
	CreateTime   string           `json:"create_time"`
	Private      bool             `json:"private"`
	Relations    FirewallRelation `json:"relations"`
	Description  string           `json:"description"`
	LocationName string           `json:"location_name"`
	Name         string           `json:"name"`
}

//FirewallRules is JSON struct of a list of firewall's rules
type FirewallRules struct {
	RulesV6In  []FirewallRuleProperties `json:"rules-v6-in,omitempty"`
	RulesV6Out []FirewallRuleProperties `json:"rules-v6-out,omitempty"`
	RulesV4In  []FirewallRuleProperties `json:"rules-v4-in,omitempty"`
	RulesV4Out []FirewallRuleProperties `json:"rules-v4-out,omitempty"`
}

//FirewallRuleProperties is JSON struct of a firewall's rule properties
type FirewallRuleProperties struct {
	Protocol string `json:"protocol,omitempty"`
	DstPort  string `json:"dst_port,omitempty"`
	SrcPort  string `json:"src_port,omitempty"`
	SrcCidr  string `json:"src_cidr,omitempty"`
	Action   string `json:"action"`
	Comment  string `json:"comment,omitempty"`
	DstCidr  string `json:"dst_cidr,omitempty"`
	Order    int    `json:"order"`
}

//FirewallRelation is a JSON struct of a list of firewall's relations
type FirewallRelation struct {
	Networks []NetworkInFirewall `json:"networks"`
}

//NetworkInFirewall is a JSON struct of a firewall's relation
type NetworkInFirewall struct {
	CreateTime  string `json:"create_time"`
	NetworkUUID string `json:"network_uuid"`
	NetworkName string `json:"network_name"`
	ObjectUUID  string `json:"object_uuid"`
	ObjectName  string `json:"object_name"`
}

//FirewallCreateRequest is JSON struct of a request for creating a firewall
type FirewallCreateRequest struct {
	Name   string        `json:"name"`
	Labels []string      `json:"labels,omitempty"`
	Rules  FirewallRules `json:"rules"`
}

//FirewallCreateResponse is JSON struct of a response for creating a firewall
type FirewallCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
}

//FirewallUpdateRequest is JSON struct of a request for updating a firewall
type FirewallUpdateRequest struct {
	Name   string        `json:"name,omitempty"`
	Labels []string      `json:"labels,omitempty"`
	Rules  FirewallRules `json:"rules,omitempty"`
}

//GetFirewallList gets a list of available firewalls
func (c *Client) GetFirewallList() ([]Firewall, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase),
		method: http.MethodGet,
	}
	var response FirewallList
	var firewalls []Firewall
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		firewalls = append(firewalls, Firewall{Properties: properties})
	}
	return firewalls, err
}

//GetFirewall gets a specific firewall based on given id
func (c *Client) GetFirewall(id string) (Firewall, error) {
	if !isValidUUID(id) {
		return Firewall{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodGet,
	}
	var response Firewall
	err := r.execute(*c, &response)
	return response, err
}

//CreateFirewall creates a new firewall
func (c *Client) CreateFirewall(body FirewallCreateRequest) (FirewallCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiFirewallBase),
		method: http.MethodPost,
		body:   body,
	}
	var response FirewallCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return FirewallCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//UpdateFirewall update a specific firewall
func (c *Client) UpdateFirewall(id string, body FirewallUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteFirewall delete a specific firewall
func (c *Client) DeleteFirewall(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetFirewallEventList get list of a firewall's events
func (c *Client) GetFirewallEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiFirewallBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var firewallEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		firewallEvents = append(firewallEvents, Event{Properties: properties})
	}
	return firewallEvents, err
}
