package gsclient

import (
	"log"
)

type Sshkey struct {
	Properties SshkeyProperties `json:"sshkey"`
}

type SshkeyProperties struct {
	Name       string `json:"name"`
	ObjectUuid string `json:"object_uuid"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
	ChangeTime string `json:"change_time"`
	Sshkey     string `json:"sshkey"`
}

func (c *Client) GetSshkey(id string) (*Sshkey, error) {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Sshkey)
	err := r.execute(*c, &response)

	log.Printf("Received sshkey: %v", response)

	return response, err
}

func (c *Client) CreateSshkey(body map[string]interface{}) (*CreateResponse, error) {
	r := Request{
		uri:    "/objects/sshkeys",
		method: "POST",
		body:   body,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *Client) DeleteSshkey(id string) error {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateSshkey(id string, body map[string]interface{}) error {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}
