package gsclient

import (
	"fmt"
)

type Servers struct {
	List map[string]ServerProperties `json:"servers"`
}

type Server struct {
	Properties ServerProperties `json:"server"`
}

type ServerProperties struct {
	ObjectUuid           string          `json:"object_uuid"`
	Name                 string          `json:"name"`
	Memory               int             `json:"memory"`
	Cores                int             `json:"cores"`
	HardwareProfile      string          `json:"hardware_profile"`
	Status               string          `json:"status"`
	LocationUuid         string          `json:"location_uuid"`
	Power                bool            `json:"power"`
	CurrentPrice         float64         `json:"current_price"`
	AvailablityZone      string          `json:"availability_zone"`
	AutoRecovery         bool            `json:"auto_recovery"`
	Legacy               bool            `json:"legacy"`
	ConsoleToken         string          `json:"console_token"`
	UsageInMinutesMemory int             `json:"usage_in_minutes_memory"`
	UsageInMinutesCores  int             `json:"usage_in_minutes_cores"`
	Labels               []string        `json:"labels"`
	Relations            ServerRelations `json:"relations"`
}

type ServerRelations struct {
	IsoImages []ServerIsoImage `json:"isoimages"`
	Networks  []ServerNetwork  `json:"networks"`
	PublicIps []ServerIp       `json:"public_ips"`
	Storages  []ServerStorage  `json:"storages"`
}

type ServerStorage struct {
	ObjectUuid       string `json:"object_uuid"`
	ObjectName       string `json:"object_name"`
	Capacity         int    `json:"capacity"`
	StorageType      string `json:"storage_type"`
	Target           int    `json:"target"`
	Lun              int    `json:"lun"`
	Controller       int    `json:"controller"`
	CreateTime       string `json:"create_time"`
	BootDevice       bool   `json:"bootdevice"`
	Bus              int    `json:"bus"`
	LastUsedTemplate string `json:"last_used_template"`
	LicenseProductNo int    `json:"license_product_no"`
	ServerUuid       string `json:"server_uuid"`
}

type ServerIsoImage struct {
	ObjectUuid string `json:"object_uuid"`
	ObjectName string `json:"object_name"`
	Private    bool   `json:"private"`
	CreateTime string `json:"create_time"`
}

type ServerNetwork struct {
	L2security           bool   `json:"l2security"`
	ServerUuid           string `json:"server_uuid"`
	CreateTime           string `json:"create_time"`
	PublicNet            bool   `json:"public_net"`
	FirewallTemplateUuid string `json:"firewall_template_uuid,omitempty"`
	ObjectName           string `json:"object_name"`
	Mac                  string `json:"mac"`
	BootDevice           bool   `json:"bootdevice"`
	PartnerUuid          string `json:"partner_uuid"`
	Ordering             int    `json:"ordering"`
	Firewall             string `json:"firewall,omitempty"`
	NetworkType          string `json:"network_type"`
	NetworkUuid          string `json:"network_uuid"`
	ObjectUuid           string `json:"object_uuid"`
	//L3security           []interface{} `json:"l3security"`
	//Vlan                 int          `json:"vlan,omitempty"`
	//Vxlan                int          `json:"vxlan,omitempty"`
	//Mcast                string       `json:"mcast, omitempty"`
}

type ServerIp struct {
	ServerUuid string `json:"server_uuid"`
	CreateTime string `json:"create_time"`
	Prefix     string `json:"prefix"`
	Family     int    `json:"family"`
	ObjectUuid string `json:"object_uuid"`
	Ip         string `json:"ip"`
}

type ServerCreateRequest struct {
	Name            string                       `json:"name"`
	Memory          int                          `json:"memory"`
	Cores           int                          `json:"cores"`
	LocationUuid    string                       `json:"location_uuid"`
	HardwareProfile string                       `json:"hardware_profile,omitempty"`
	AvailablityZone string                       `json:"availability_zone,omitempty"`
	Labels          []interface{}                `json:"labels,omitempty"`
	Relations       ServerCreateRequestRelations `json:"relations,omitempty"`
}

type ServerCreateRequestRelations struct {
	IsoImages []ServerCreateRequestIsoimage `json:"isoimages"`
	Networks  []ServerCreateRequestNetwork  `json:"networks"`
	PublicIps []ServerCreateRequestIp       `json:"public_ips"`
	Storages  []ServerCreateRequestStorage  `json:"storages"`
}

type ServerCreateRequestStorage struct {
	StorageUuid string `json:"storage_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

type ServerCreateRequestNetwork struct {
	NetworkUuid string `json:"network_uuid,omitempty"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

type ServerCreateRequestIp struct {
	IpaddrUuid string `json:"ipaddr_uuid,omitempty"`
}

type ServerCreateRequestIsoimage struct {
	IsoimageUuid string `json:"isoimage_uuid,omitempty"`
}

type ServerUpdateRequest struct {
	Name            string        `json:"name,omitempty"`
	AvailablityZone string        `json:"availability_zone,omitempty"`
	Memory          int           `json:"memory,omitempty"`
	Cores           int           `json:"cores,omitempty"`
	Labels          []interface{} `json:"labels"`
}

func (c *Client) GetServer(id string) (*Server, error) {
	if id == "" {
		return nil, fmt.Errorf(
			"Can't read without id", nil)
	}
	r := Request{
		uri:    apiServerBase + "/" + id,
		method: "GET",
	}

	response := new(Server)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) GetServerList() ([]Server, error) {
	r := Request{
		uri:    apiServerBase,
		method: "GET",
	}

	response := new(Servers)
	err := r.execute(*c, &response)

	list := []Server{}
	for _, properties := range response.List {
		server := Server{
			Properties: properties,
		}
		list = append(list, server)
	}

	return list, err
}

func (c *Client) CreateServer(body ServerCreateRequest) (*CreateResponse, error) {
	r := Request{
		uri:    apiServerBase,
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

func (c *Client) DeleteServer(id string) error {
	r := Request{
		uri:    apiServerBase + "/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateServer(id string, body ServerUpdateRequest) error {
	r := Request{
		uri:    apiServerBase + "/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) StopServer(id string) error {
	server, err := c.GetServer(id)
	if err != nil {
		return err
	}
	if !server.Properties.Power {
		return nil
	}

	body := map[string]interface{}{
		"power": false,
	}
	r := Request{
		uri:    apiServerBase + "/" + id + "/power",
		method: "PATCH",
		body:   body,
	}

	err = r.execute(*c, nil)
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
		uri:    apiServerBase + "/" + id + "/shutdown",
		method: "PATCH",
		body:   new(map[string]string),
	}

	err = r.execute(*c, nil)
	if err != nil {
		return err
	}

	//If we get an error, which includes a timeout, power off the server instead
	err = c.WaitForServerPowerStatus(id, false)
	if err != nil {
		return c.StopServer(id)
	}

	return nil
}

func (c *Client) StartServer(id string) error {
	server, err := c.GetServer(id)
	if err != nil {
		return err
	}
	if server.Properties.Power {
		return nil
	}

	body := map[string]interface{}{
		"power": true,
	}
	r := Request{
		uri:    apiServerBase + "/" + id + "/power",
		method: "PATCH",
		body:   body,
	}

	err = r.execute(*c, nil)
	if err != nil {
		return err
	}

	return c.WaitForServerPowerStatus(id, true)
}

func (c *Client) LinkStorage(serverid string, storageid string, bootdevice bool) error {
	body := map[string]interface{}{
		"object_uuid": storageid,
		"bootdevice":  bootdevice,
	}
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/storages",
		method: "POST",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) UnlinkStorage(serverid string, storageid string) error {
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/storages/" + storageid,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) LinkNetwork(serverid string, networkid string, bootdevice bool) error {
	body := map[string]interface{}{
		"object_uuid": networkid,
		"bootdevice":  bootdevice,
	}
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/networks",
		method: "POST",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) UnlinkNetwork(serverid string, networkid string) error {
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/networks/" + networkid,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) LinkIsoimage(serverid string, isoimageid string) error {
	body := map[string]interface{}{
		"object_uuid": isoimageid,
	}
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/isoimages",
		method: "POST",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) UnlinkIsoimage(serverid string, isoimageid string) error {
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/isoimages/" + isoimageid,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) LinkIp(serverid string, ipid string) error {
	body := map[string]interface{}{
		"object_uuid": ipid,
	}
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/ips",
		method: "POST",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) UnlinkIp(serverid string, ipid string) error {
	r := Request{
		uri:    apiServerBase + "/" + serverid + "/ips/" + ipid,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}
