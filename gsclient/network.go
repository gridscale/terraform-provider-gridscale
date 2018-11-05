package gsclient

import (
	"log"
)

type Network struct {
	Body struct {
		ObjectUuid   string `json:"object_uuid"`
		Name         string `json:"name"`
		LocationUuid string `json:"location_uuid"`
	} `json:"network"`
}

func (c *Client) ReadNetwork(id string) (*Network, error) {
	r := Request{
		uri:    "/objects/networks/" + id,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Network)
	err := r.execute(*c, &response)

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

	return response, err
}

func (c *Client) DestroyNetwork(id string) error {
	r := Request{
		uri:    "/objects/networks/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}
