package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ServerIPRelationList JSON struct of a list of relations between a server and IP addresses
type ServerIPRelationList struct {
	List []ServerIPRelationProperties `json:"ip_relations"`
}

//ServerIPRelation JSON struct of a single relation between a server and a IP address
type ServerIPRelation struct {
	Properties ServerIPRelationProperties `json:"ip_relation"`
}

//ServerIPRelationProperties JSON struct of properties of a relation between a server and a IP address
type ServerIPRelationProperties struct {
	ServerUUID string `json:"server_uuid"`
	CreateTime string `json:"create_time"`
	Prefix     string `json:"prefix"`
	Family     int    `json:"family"`
	ObjectUUID string `json:"object_uuid"`
	IP         string `json:"ip"`
}

//ServerIPRelationCreateRequest JSON struct of request for creating a relation between a server and a IP address
type ServerIPRelationCreateRequest struct {
	ObjectUUID string `json:"object_uuid"`
}

//GetServerIPList gets a list of a specific server's IPs
func (c *Client) GetServerIPList(id string) ([]ServerIPRelationProperties, error) {
	if id == "" {
		return nil, errors.New("'id' is required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodGet,
	}
	var response ServerIPRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerIP gets an IP of a specific server
func (c *Client) GetServerIP(serverID, ipID string) (ServerIPRelationProperties, error) {
	if serverID == "" || ipID == "" {
		return ServerIPRelationProperties{}, errors.New("'serverID' and 'ipID' are required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "ips", ipID),
		method: http.MethodGet,
	}
	var response ServerIPRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//CreateServerIP create a link between a server and an IP
func (c *Client) CreateServerIP(id string, body ServerIPRelationCreateRequest) error {
	if id == "" || body.ObjectUUID == "" {
		return errors.New("'server_id' and 'ip_id' are required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "ips"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerIP delete a link between a server and an IP
func (c *Client) DeleteServerIP(serverID, ipID string) error {
	if serverID == "" || ipID == "" {
		return errors.New("'serverID' and 'ipID' are required")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "ips", ipID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkIP attaches an IP to a server
func (c *Client) LinkIP(serverID string, ipID string) error {
	body := ServerIPRelationCreateRequest{
		ObjectUUID: ipID,
	}
	return c.CreateServerIP(serverID, body)
}

//UnlinkIP removes a link between an IP and a server
func (c *Client) UnlinkIP(serverID string, ipID string) error {
	return c.DeleteServerIP(serverID, ipID)
}
