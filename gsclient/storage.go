package gsclient

import (
	"log"
)
type Storage struct {
	Body struct {
		ObjectUuid		string	`json:"object_uuid"`
		Name			string	`json:"name"`
		Capacity		string	`json:"capacity"`
		LocationUuid	string	`json:"location_uuid"`
	} `json:"storage"`
}

func (c *Client) ReadStorage(id string) (*Storage, error) {
	r := Request{
		uri: 			"/objects/storages/" + id,
		method: 		"GET",
	}
	log.Printf("%v", r)

	response := new(Storage)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) CreateStorage(body map[string]interface{}) (*CreateResponse, error) {
	r := Request{
		uri: 			"/objects/storages",
		method: 		"POST",
		body:			body,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) DestroyStorage(id string) error {
	r := Request{
		uri: 			"/objects/storages/" + id,
		method: 		"DELETE",
	}

	return r.execute(*c, nil)
}