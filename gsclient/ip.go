package gsclient

import (
	"log"
)

type Ip struct {
	Properties IpProperties `json:"ip"`
}

type IpProperties struct {
	LocationCountry string `json:"location_country"`
	LocationUuid    string `json:"location_uuid"`
	ObjectUuid      string `json:"object_uuid"`
	ReverseDns      string `json:"reverse_dns"`
	Family          int    `json:"family"`
	Status          string `json:"status"`
	CreateTime      string `json:"create_time"`
	Failover        bool   `json:"failover"`
	ChangeTime      string `json:"change_time"`
	LocationIata    string `json:"location_iata"`
	LocationName    string `json:"location_name"`
	Prefix          string `json:"prefix"`
	Ip              string `json:"ip"`
}

type IpCreateResponse struct {
	ObjectUuid string `json:"object_uuid"`
	Prefix     string `json:"prefix"`
	Ip         string `json:"ip"`
}

func (c *Client) GetIp(id string) (*Ip, error) {
	r := Request{
		uri:    "/objects/ips/" + id,
		method: "GET",
	}
	log.Printf("%v", r)

	response := new(Ip)
	err := r.execute(*c, &response)

	log.Printf("Received network: %v", response)

	return response, err
}

func (c *Client) CreateIp(body map[string]interface{}) (*IpCreateResponse, error) {
	r := Request{
		uri:    "/objects/ips",
		method: "POST",
		body:   body,
	}

	response := new(IpCreateResponse)
	err := r.execute(*c, &response)
	if err != nil {
		return nil, err
	}

	return response, err
}

func (c *Client) DeleteIp(id string) error {
	r := Request{
		uri:    "/objects/ips/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateIp(id string, body map[string]interface{}) error {
	r := Request{
		uri:    "/objects/ips/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}
