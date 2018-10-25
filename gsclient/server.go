package gsclient

import (
	"fmt"
	"log"
)

type Server struct {
	Properties			ServerProperties	`json:"server"`
}

type ServerProperties struct {
	ObjectUuid		string 				`json:"object_uuid"`
	Name			string 				`json:"name"`
	Memory			int    				`json:"memory"`
	Cores			int    				`json:"cores"`
	HardwareProfile	string				`json:"hardware_profile"`
	Status			string				`json:"status"`
	LocationUuid	string				`json:"location_uuid"`
	Power			bool				`json:"power"`
	CurrentPrice	float32				`json:"current_price"`
	Relations		ServerRelations		`json:"relations"`
}

type ServerRelations struct {
	IsoImages	[]interface{}	`json:"isoimages"`
	Networks	[]interface{}	`json:"networks"`
	PublicIps	[]interface{}	`json:"public_ips"`
	Storages	[]ServerStorage	`json:"storages"`
}

type ServerStorage struct {
	StorageUuid	string	`json:"storage_uuid"`
	BootDevice	bool	`json:"bootdevice"`
}

type ServerCreateRequest struct {
	Name			string 				`json:"name"`
	Memory			int    				`json:"memory"`
	Cores			int    				`json:"cores"`
	LocationUuid	string				`json:"location_uuid"`
	Relations		ServerRelations		`json:"relations"`
}

func (c *Client) GetServer(id string) (*Server, error) {
	if id == "" {
		return nil, fmt.Errorf(
			"Can't read without id", nil)
	}
	r := Request{
		uri: 			"/objects/servers/" + id,
		method: 		"GET",
	}
	log.Printf("%v", r)

	response := new(Server)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) CreateServer(s ServerCreateRequest) (*CreateResponse, error) {
	r := Request{
		uri: 			"/objects/servers",
		method: 		"POST",
		body:			s,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)

	c.WaitForState(response.ServerUuid, c.serverIsProvisioned)

	return response, err
}

func (c *Client) DeleteServer(id string) error {
	r := Request{
		uri: 			"/objects/servers/" + id,
		method: 		"DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateServer(s *Server) (*Server, error) {
	body := map[string]interface{}{}


	r := Request{
		uri:			"/objects/servers/" + s.Properties.ObjectUuid,
		method:			"PATCH",
		body:			body,
	}

	return s, r.execute(*c, s)
}

func (c *Client) StopServer(s Server) error {
	if !s.Properties.Power{
		return nil
	}

	body := map[string]interface{}{
		"power":	false,
	}
	r := Request{
		uri:			"/objects/servers/" + s.Properties.ObjectUuid + "/power",
		method:			"PATCH",
		body:			body,
	}

	return r.execute(*c, nil)
}

func (c *Client) StartServer(s Server) error {
	if s.Properties.Power{
		return nil
	}

	body := map[string]interface{}{
		"power":	true,
	}
	r := Request{
		uri:			"/objects/servers/" + s.Properties.ObjectUuid + "/power",
		method:			"PATCH",
		body:			body,
	}

	return r.execute(*c, nil)
}

func (c *Client) serverIsProvisioned(id string, inProgress chan bool) {
	server, _ := c.GetServer(id)

	if server.Properties.Status != "active" {
		inProgress <- false
	}
}

func (c *Client) serverIsDeletable(id string, inProgress chan bool) {
	server, _ := c.GetServer(id)

	if server.Properties.Power == false {
		inProgress <- false
	}
}