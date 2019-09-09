package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ServerNetworkRelationList JSON struct of a list of relations between a server and networks
type ServerNetworkRelationList struct {
	List []ServerNetworkRelationProperties `json:"network_relations"`
}

//ServerNetworkRelation JSON struct of a single relation between a server and a network
type ServerNetworkRelation struct {
	Properties ServerNetworkRelationProperties `json:"network_relation"`
}

//ServerNetworkRelationProperties JSON struct of properties of a relation between a server and a network
type ServerNetworkRelationProperties struct {
	L2security           bool     `json:"l2security"`
	ServerUUID           string   `json:"server_uuid"`
	CreateTime           string   `json:"create_time"`
	PublicNet            bool     `json:"public_net"`
	FirewallTemplateUUID string   `json:"firewall_template_uuid,omitempty"`
	ObjectName           string   `json:"object_name"`
	Mac                  string   `json:"mac"`
	BootDevice           bool     `json:"bootdevice"`
	PartnerUUID          string   `json:"partner_uuid"`
	Ordering             int      `json:"ordering"`
	Firewall             string   `json:"firewall,omitempty"`
	NetworkType          string   `json:"network_type"`
	NetworkUUID          string   `json:"network_uuid"`
	ObjectUUID           string   `json:"object_uuid"`
	L3security           []string `json:"l3security"`
	//Vlan                 int          `json:"vlan,omitempty"`
	//Vxlan                int          `json:"vxlan,omitempty"`
	//Mcast                string       `json:"mcast, omitempty"`
}

//ServerNetworkRelationCreateRequest JSON struct of a request for creating a relation between a server and a network
type ServerNetworkRelationCreateRequest struct {
	ObjectUUID           string        `json:"object_uuid"`
	Ordering             int           `json:"ordering,omitempty"`
	BootDevice           bool          `json:"bootdevice,omitempty"`
	L3security           []string      `json:"l3security,omitempty"`
	Firewall             FirewallRules `json:"firewall,omitempty"`
	FirewallTemplateUUID string        `json:"firewall_template_uuid,omitempty"`
}

//ServerNetworkRelationUpdateRequest JSON struct of a request for updating a relation between a server and a network
type ServerNetworkRelationUpdateRequest struct {
	Ordering             int           `json:"ordering"`
	BootDevice           bool          `json:"bootdevice"`
	L3security           []string      `json:"l3security"`
	Firewall             FirewallRules `json:"firewall"`
	FirewallTemplateUUID string        `json:"firewall_template_uuid"`
}

//GetServerNetworkList gets a list of a specific server's networks
func (c *Client) GetServerNetworkList(id string) ([]ServerNetworkRelationProperties, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodGet,
	}
	var response ServerNetworkRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerNetwork gets a network of a specific server
func (c *Client) GetServerNetwork(serverID, networkID string) (ServerNetworkRelationProperties, error) {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return ServerNetworkRelationProperties{}, errors.New("'serverID' or 'networksID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "networks", networkID),
		method: http.MethodGet,
	}
	var response ServerNetworkRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerNetwork updates a link between a network and a server
func (c *Client) UpdateServerNetwork(serverID, networkID string, body ServerNetworkRelationUpdateRequest) error {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return errors.New("'serverID' or 'networksID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "networks", networkID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerNetwork creates a link between a network and a storage
func (c *Client) CreateServerNetwork(id string, body ServerNetworkRelationCreateRequest) error {
	if !isValidUUID(id) || !isValidUUID(body.ObjectUUID) {
		return errors.New("'serverID' or 'network_id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "networks"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerNetwork deletes a link between a network and a server
func (c *Client) DeleteServerNetwork(serverID, networkID string) error {
	if !isValidUUID(serverID) || !isValidUUID(networkID) {
		return errors.New("'serverID' or 'networkID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "networks", networkID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkNetwork attaches a network to a server
func (c *Client) LinkNetwork(serverID, networkID, firewallTemplate string, bootdevice bool, order int,
	l3security []string, firewall FirewallRules) error {
	body := ServerNetworkRelationCreateRequest{
		ObjectUUID:           networkID,
		Ordering:             order,
		BootDevice:           bootdevice,
		L3security:           l3security,
		FirewallTemplateUUID: firewallTemplate,
		Firewall:             firewall,
	}
	return c.CreateServerNetwork(serverID, body)
}

//UnlinkNetwork removes the link between a network and a server
func (c *Client) UnlinkNetwork(serverID string, networkID string) error {
	return c.DeleteServerNetwork(serverID, networkID)
}
