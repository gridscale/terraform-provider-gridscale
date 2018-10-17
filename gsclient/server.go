package gsclient

import (
	"log"
)
type Server struct {
	Body struct {
		ObjectUuid string `json:"object_uuid"`
		Name       string `json:"name"`
		Memory     int    `json:"memory"`
		Cores      int    `json:"cores"`
		HardwareProfile	string	`json:"hardware_prfile"`
	} `json:"server"`
}

type CreateResponse struct {
	ObjectUuid  string `json:"object_uuid"`
	RequestUuid string `json:"request_uuid"`
}

func (c *Client) ReadServer(id string) (*Server, error) {
	r := Request{
		uri: 			"/objects/servers/" + id,
		method: 		"GET",
	}
	log.Printf("%v", r)

	response := new(Server)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) CreateServer(body map[string]interface{}) (string, error) {
	r := Request{
		uri: 			"/objects/servers",
		method: 		"POST",
		body:			body,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)

	return response.ObjectUuid, err
}

func (c *Client) DestroyServer(id string) error {
	r := Request{
		uri: 			"/objects/servers/" + id,
		method: 		"DELETE",
	}

	return r.execute(*c, nil)
}