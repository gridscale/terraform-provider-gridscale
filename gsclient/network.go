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
	LocationCountry string `json:"location_country"`
	LocationUuid    string `json:"location_uuid"`
	PublicNet       bool   `json:"public_net"`
	ObjectUuid      string `json:"object_uuid"`
	NetworkType     string `json:"network_type"`
	Name            string `json:"name"`
	Status          string `json:"status"`
	CreateTime      string `json:"create_time"`
	L2Security      bool   `json:"l2security"`
	ChangeTime      string `json:"change_time"`
	LocationIata    string `json:"location_iata"`
	LocationName    string `json:"location_name"`
	DeleteBlock     bool   `json:"delete_block"`
}

func (c *Client) GetNetwork(id string) (*Network, error) {
	r := Request{
		uri:    "/objects/networks/" + id,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Network)
	err := r.execute(*c, &response)

	log.Printf("Received network: %v", response)

	return response, err
}

func (c *Client) CreateNetwork(body map[string]interface{}) (*CreateResponse, error) {
	r := Request{
		uri:    "/objects/networks",
		method: "POST",
		body:   body,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)
	if err != nil {
		return nil, err
	}

	err = c.WaitForRequestCompletion(*response)

	return response, err
}

func (c *Client) DeleteNetwork(id string) error {
	r := Request{
		uri:    "/objects/networks/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateNetwork(id string, body map[string]interface{}) error {
	r := Request{
		uri:    "/objects/networks/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) GetNetworkList() (*Networks, error) {
	r := Request{
		uri:    "/objects/networks/",
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Networks)
	err := r.execute(*c, &response)

	log.Printf("Received networks: %v", response)

	return response, err
}

func (c *Client) GetNetworkPublic() (string, error) {
	networks, err := c.GetNetworkList()
	if err != nil {
		return "", err
	}
	for _, value := range networks.List {
		if value.Name == "Public Network" {
			return value.ObjectUuid, nil
		}
	}

	return "", fmt.Errorf("Public Network not found")
}