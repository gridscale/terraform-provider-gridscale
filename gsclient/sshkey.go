package gsclient

import (
	"log"
)

type Sshkeys struct {
	List map[string]SshkeyProperties `json:"sshkeys"`
}

type Sshkey struct {
	Properties SshkeyProperties `json:"sshkey"`
}

type SshkeyProperties struct {
	Name       string   `json:"name"`
	ObjectUuid string   `json:"object_uuid"`
	Status     string   `json:"status"`
	CreateTime string   `json:"create_time"`
	ChangeTime string   `json:"change_time"`
	Sshkey     string   `json:"sshkey"`
	Labels     []string `json:"labels"`
}

type SshkeyCreateRequest struct {
	Name   string        `json:"name"`
	Sshkey string        `json:"sshkey"`
	Labels []interface{} `json:"labels,omitempty"`
}

type SshkeyUpdateRequest struct {
	Name   string        `json:"name,omitempty"`
	Sshkey string        `json:"sshkey,omitempty"`
	Labels []interface{} `json:"labels,omitempty"`
}

func (c *Client) GetSshkey(id string) (*Sshkey, error) {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "GET",
	}

	response := new(Sshkey)
	err := r.execute(*c, &response)

	log.Printf("Received sshkey: %v", response)

	return response, err
}

func (c *Client) GetSshkeyList() (*Sshkeys, error) {
	r := Request{
		uri:    "/objects/sshkeys/",
		method: "GET",
	}

	response := new(Sshkeys)
	err := r.execute(*c, &response)

	log.Printf("Received sshkeys: %v", response)

	return response, err
}

func (c *Client) CreateSshkey(body SshkeyCreateRequest) (*CreateResponse, error) {
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

	err = c.WaitForRequestCompletion(response.RequestUuid)

	return response, err
}

func (c *Client) DeleteSshkey(id string) error {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateSshkey(id string, body SshkeyUpdateRequest) error {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}
