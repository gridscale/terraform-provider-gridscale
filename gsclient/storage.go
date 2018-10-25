package gsclient

import (
	"log"
)
type Storage struct {
	Properties		StorageProperties	`json:"storage"`
}

type StorageProperties struct {
	ObjectUuid		string	`json:"object_uuid"`
	Name			string	`json:"name"`
	Capacity		string	`json:"capacity"`
	LocationUuid	string	`json:"location_uuid"`
	Status			string	`json:"status"`
}

func (c *Client) GetStorage(id string) (*Storage, error) {
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

	c.WaitForState(response.ObjectUuid, c.storageIsProvisioned)

	return response, err
}

func (c *Client) DeleteStorage(id string) error {
	r := Request{
		uri: 			"/objects/storages/" + id,
		method: 		"DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) storageIsProvisioned(id string, inProgress chan bool) {
	storage, _ := c.GetStorage(id)

	if storage.Properties.Status != "active" {
		inProgress <- false
	}
}