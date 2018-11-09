package gsclient

import (
	"fmt"
	"log"
)

type Servers struct {
	List map[string]ServerProperties `json:"servers"`
}

type Server struct {
	Properties ServerProperties `json:"server"`
}

type ServerProperties struct {
	ObjectUuid      string          `json:"object_uuid"`
	Name            string          `json:"name"`
	Memory          int             `json:"memory"`
	Cores           int             `json:"cores"`
	HardwareProfile string          `json:"hardware_profile"`
	Status          string          `json:"status"`
	LocationUuid    string          `json:"location_uuid"`
	Power           bool            `json:"power"`
	CurrentPrice    float64         `json:"current_price"`
	Relations       ServerRelations `json:"relations"`
}

type ServerRelations struct {
	IsoImages []ServerIsoImage `json:"isoimages"`
	Networks  []ServerNetwork  `json:"networks"`
	PublicIps []ServerIp       `json:"public_ips"`
	Storages  []ServerStorage  `json:"storages"`
}

type ServerStorage struct {
	StorageUuid string `json:"storage_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

type ServerIsoImage struct {
	IsoImageUuid string `json:"isoimage_uuid,omitempty"`
}

type ServerNetwork struct {
	NetworkUuid string `json:"network_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

type ServerIp struct {
	IpaddrUuid string `json:"ipaddr_uuid,omitempty"`
}

type ServerCreateRequest struct {
	Name            string          `json:"name"`
	Memory          int             `json:"memory"`
	Cores           int             `json:"cores"`
	LocationUuid    string          `json:"location_uuid"`
	HardwareProfile string          `json:"hardware_profile"`
	Relations       ServerRelations `json:"relations,omitempty"`
}

func (c *Client) GetServer(id string) (*Server, error) {
	if id == "" {
		return nil, fmt.Errorf(
			"Can't read without id", nil)
	}
	r := Request{
		uri:    "/objects/servers/" + id,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Server)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) GetServerList() (*Servers, error) {
	r := Request{
		uri:    "/objects/servers/",
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Servers)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) CreateServer(s ServerCreateRequest) (*CreateResponse, error) {
	r := Request{
		uri:    "/objects/servers",
		method: "POST",
		body:   s,
	}

	response := new(CreateResponse)
	err := r.execute(*c, &response)
	if err != nil {
		return nil, err
	}

	err = c.WaitForRequestCompletion(response.RequestUuid)

	return response, err
}

func (c *Client) DeleteServer(id string) error {
	r := Request{
		uri:    "/objects/servers/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateServer(id string, body map[string]interface{}) error {
	r := Request{
		uri:    "/objects/servers/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) StopServer(id string) error {
	//Make sure the server exists and that it isn't already in the state we need it to be
	body := map[string]interface{}{
		"power": false,
	}
	r := Request{
		uri:    "/objects/servers/" + id + "/power",
		method: "PATCH",
		body:   body,
	}

	err := r.execute(*c, nil)
	if err != nil {
		return err
	}

	return c.WaitForServerPowerStatus(id, false)
}

func (c *Client) ShutdownServer(id string) error {
	//Make sure the server exists and that it isn't already in the state we need it to be
	server, err := c.GetServer(id)
	if err != nil {
		return err
	}
	if !server.Properties.Power {
		return nil
	}

	r := Request{
		uri:    "/objects/servers/" + id + "/shutdown",
		method: "PATCH",
	}

	err = r.execute(*c, nil)
	if err != nil {
		return err
	}

	//If we get an error, which includes a timeout, power off the server instead
	err = c.WaitForServerPowerStatus(id, false)
	if err != nil {
		c.StopServer(id)
	}

	return nil
}

func (c *Client) StartServer(id string) error {
	body := map[string]interface{}{
		"power": true,
	}
	r := Request{
		uri:    "/objects/servers/" + id + "/power",
		method: "PATCH",
		body:   body,
	}

	err := r.execute(*c, nil)
	if err != nil {
		return err
	}

	return c.WaitForServerPowerStatus(id, true)
}
