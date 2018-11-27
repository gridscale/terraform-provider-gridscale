package gsclient

import (
	"fmt"
	"log"
)

type Networks struct {
	List map[string]NetworkProperties `json:"networks"`
}

type Network struct {
	Properties NetworkProperties `json:"network"`
}

type NetworkProperties struct {
	LocationCountry string           `json:"location_country"`
	LocationUuid    string           `json:"location_uuid"`
	PublicNet       bool             `json:"public_net"`
	ObjectUuid      string           `json:"object_uuid"`
	NetworkType     string           `json:"network_type"`
	Name            string           `json:"name"`
	Status          string           `json:"status"`
	CreateTime      string           `json:"create_time"`
	L2Security      bool             `json:"l2security"`
	ChangeTime      string           `json:"change_time"`
	LocationIata    string           `json:"location_iata"`
	LocationName    string           `json:"location_name"`
	DeleteBlock     bool             `json:"delete_block"`
	Labels          []string         `json:"labels"`
	Relations       NetworkRelations `json:"relations"`
}

type NetworkRelations struct {
	Vlans   []NetworkVlan   `json:"vlans"`
	Servers []NetworkServer `json:"servers"`
}

type NetworkVlan struct {
	Vlan       int    `json:"vlan"`
	TenantName string `json:"tenant_name"`
	TenantUuid string `json:"tenant_uuid"`
}

type NetworkServer struct {
	ObjectUuid  string   `json:"object_uuid"`
	Mac         string   `json:"mac"`
	Bootdevice  bool     `json:"bootdevice"`
	CreateTime  string   `json:"create_time"`
	L3security  []string `json:"l3security"`
	ObjectName  string   `json:"object_name"`
	NetworkUuid string   `json:"network_uuid"`
	Ordering    int      `json:"ordering"`
}

type NetworkCreateRequest struct {
	Name         string        `json:"name"`
	Labels       []interface{} `json:"labels,omitempty"`
	LocationUuid string        `json:"location_uuid"`
	L2Security   bool          `json:"l2security,omitempty"`
}

type NetworkUpdateRequest struct {
	Name       string        `json:"name,omitempty"`
	Labels     []interface{} `json:"labels"`
	L2Security bool          `json:"l2security"`
}

func (c *Client) GetNetwork(id string) (*Network, error) {
	r := Request{
		uri:    apiNetworkBase + "/" + id,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Network)
	err := r.execute(*c, &response)

	log.Printf("Received network: %v", response)

	return response, err
}

func (c *Client) CreateNetwork(body NetworkCreateRequest) (*CreateResponse, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: "POST",
		body:   body,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)
	if err != nil {
		return nil, err
	}

	err = c.WaitForRequestCompletion(response.RequestUuid)

	return response, err
}

func (c *Client) DeleteNetwork(id string) error {
	r := Request{
		uri:    apiNetworkBase + "/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateNetwork(id string, body NetworkUpdateRequest) error {
	r := Request{
		uri:    apiNetworkBase + "/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) GetNetworkList() ([]Network, error) {
	r := Request{
		uri:    apiNetworkBase,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Networks)
	err := r.execute(*c, &response)

	list := []Network{}
	for _, properties := range response.List {
		network := Network{
			Properties: properties,
		}
		list = append(list, network)
	}

	return list, err
}

func (c *Client) GetNetworkPublic() (*Network, error) {
	networks, err := c.GetNetworkList()
	if err != nil {
		return nil, err
	}
	for _, network := range networks {
		if network.Properties.PublicNet {
			return c.GetNetwork(network.Properties.ObjectUuid)
		}
	}

	return nil, fmt.Errorf("Public Network not found")
}
