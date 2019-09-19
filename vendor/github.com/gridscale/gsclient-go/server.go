package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ServerList JSON struct of a list of servers
type ServerList struct {
	List map[string]ServerProperties `json:"servers"`
}

//DeletedServerList JSON struct of a list of deleted servers
type DeletedServerList struct {
	List map[string]ServerProperties `json:"deleted_servers"`
}

//Server JSON struct of a single server
type Server struct {
	Properties ServerProperties `json:"server"`
}

//ServerProperties JSON struct of properties of a server
type ServerProperties struct {
	ObjectUUID           string          `json:"object_uuid"`
	Name                 string          `json:"name"`
	Memory               int             `json:"memory"`
	Cores                int             `json:"cores"`
	HardwareProfile      string          `json:"hardware_profile"`
	Status               string          `json:"status"`
	LocationUUID         string          `json:"location_uuid"`
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

//ServerRelations JSON struct of a list of server relations
type ServerRelations struct {
	IsoImages []ServerIsoImageRelationProperties `json:"isoimages"`
	Networks  []ServerNetworkRelationProperties  `json:"networks"`
	PublicIPs []ServerIPRelationProperties       `json:"public_ips"`
	Storages  []ServerStorageRelationProperties  `json:"storages"`
}

//ServerCreateRequest JSON struct of a request for creating a server
type ServerCreateRequest struct {
	Name            string                        `json:"name"`
	Memory          int                           `json:"memory"`
	Cores           int                           `json:"cores"`
	LocationUUID    string                        `json:"location_uuid"`
	HardwareProfile string                        `json:"hardware_profile,omitempty"`
	AvailablityZone string                        `json:"availability_zone,omitempty"`
	Labels          []string                      `json:"labels,omitempty"`
	Status          string                        `json:"status,omitempty"`
	AutoRecovery    *bool                         `json:"auto_recovery,omitempty"`
	Relations       *ServerCreateRequestRelations `json:"relations,omitempty"`
}

//ServerCreateRequestRelations JSOn struct of a list of a server's relations
type ServerCreateRequestRelations struct {
	IsoImages []ServerCreateRequestIsoimage `json:"isoimages"`
	Networks  []ServerCreateRequestNetwork  `json:"networks"`
	PublicIPs []ServerCreateRequestIP       `json:"public_ips"`
	Storages  []ServerCreateRequestStorage  `json:"storages"`
}

//ServerCreateResponse JSON struct of a response for creating a server
type ServerCreateResponse struct {
	ObjectUUID   string   `json:"object_uuid"`
	RequestUUID  string   `json:"request_uuid"`
	ServerUUID   string   `json:"server_uuid"`
	NetworkUUIDs []string `json:"network_uuids"`
	StorageUUIDs []string `json:"storage_uuids"`
	IPaddrUUIDs  []string `json:"ipaddr_uuids"`
}

//ServerPowerUpdateRequest JSON struct of a request for updating server's power state
type ServerPowerUpdateRequest struct {
	Power bool `json:"power"`
}

//ServerCreateRequestStorage JSON struct of a relation between a server and a storage
type ServerCreateRequestStorage struct {
	StorageUUID string `json:"storage_uuid"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

//ServerCreateRequestNetwork JSON struct of a relation between a server and a network
type ServerCreateRequestNetwork struct {
	NetworkUUID string `json:"network_uuid"`
	BootDevice  bool   `json:"bootdevice,omitempty"`
}

//ServerCreateRequestIP JSON struct of a relation between a server and an IP address
type ServerCreateRequestIP struct {
	IPaddrUUID string `json:"ipaddr_uuid"`
}

//ServerCreateRequestIsoimage JSON struct of a relation between a server and an ISO-Image
type ServerCreateRequestIsoimage struct {
	IsoimageUUID string `json:"isoimage_uuid"`
}

//ServerUpdateRequest JSON of a request for updating a server
type ServerUpdateRequest struct {
	Name            string   `json:"name,omitempty"`
	AvailablityZone string   `json:"availability_zone,omitempty"`
	Memory          int      `json:"memory,omitempty"`
	Cores           int      `json:"cores,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	AutoRecovery    *bool    `json:"auto_recovery,omitempty"`
}

//ServerMetricList JSON struct of a list of a server's metrics
type ServerMetricList struct {
	List []ServerMetricProperties `json:"server_metrics"`
}

//ServerMetric JSON struct of a single metric of a server
type ServerMetric struct {
	Properties ServerMetricProperties `json:"server_metric"`
}

//ServerMetricProperties JSON stru
type ServerMetricProperties struct {
	BeginTime       string `json:"begin_time"`
	EndTime         string `json:"end_time"`
	PaaSServiceUUID string `json:"paas_service_uuid"`
	CoreUsage       struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"core_usage"`
	StorageSize struct {
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"storage_size"`
}

//GetServer gets a specific server based on given list
func (c *Client) GetServer(id string) (Server, error) {
	if !isValidUUID(id) {
		return Server{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodGet,
	}
	var response Server
	err := r.execute(*c, &response)
	return response, err
}

//GetServerList gets a list of available servers
func (c *Client) GetServerList() ([]Server, error) {
	r := Request{
		uri:    apiServerBase,
		method: http.MethodGet,
	}
	var response ServerList
	var servers []Server
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{
			Properties: properties,
		})
	}
	return servers, err
}

//CreateServer create a server
func (c *Client) CreateServer(body ServerCreateRequest) (ServerCreateResponse, error) {
	//check if these slices are nil
	//make them be empty slice instead of nil
	//so that JSON structure will be valid
	if body.Relations != nil && body.Relations.PublicIPs == nil {
		body.Relations.PublicIPs = make([]ServerCreateRequestIP, 0)
	}
	if body.Relations != nil && body.Relations.Networks == nil {
		body.Relations.Networks = make([]ServerCreateRequestNetwork, 0)
	}
	if body.Relations != nil && body.Relations.IsoImages == nil {
		body.Relations.IsoImages = make([]ServerCreateRequestIsoimage, 0)
	}
	if body.Relations != nil && body.Relations.Storages == nil {
		body.Relations.Storages = make([]ServerCreateRequestStorage, 0)
	}
	r := Request{
		uri:    apiServerBase,
		method: http.MethodPost,
		body:   body,
	}
	var response ServerCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return ServerCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	//this fixed the endpoint's bug temporarily when creating server with/without
	//'relations' field
	if response.ServerUUID == "" && response.ObjectUUID != "" {
		response.ServerUUID = response.ObjectUUID
	} else if response.ObjectUUID == "" && response.ServerUUID != "" {
		response.ObjectUUID = response.ServerUUID
	}
	return response, err
}

//DeleteServer deletes a specific server
func (c *Client) DeleteServer(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateServer updates a specific server
func (c *Client) UpdateServer(id string, body ServerUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetServerEventList gets a list of a specific server's events
func (c *Client) GetServerEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var serverEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		serverEvents = append(serverEvents, Event{Properties: properties})
	}
	return serverEvents, err
}

//GetServerMetricList gets a list of a specific server's metrics
func (c *Client) GetServerMetricList(id string) ([]ServerMetric, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "metrics"),
		method: http.MethodGet,
	}
	var response ServerMetricList
	var serverMetrics []ServerMetric
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		serverMetrics = append(serverMetrics, ServerMetric{Properties: properties})
	}
	return serverMetrics, err
}

//IsServerOn returns true if the server's power is on, otherwise returns false
func (c *Client) IsServerOn(id string) (bool, error) {
	server, err := c.GetServer(id)
	if err != nil {
		return false, err
	}
	return server.Properties.Power, nil
}

//setServerPowerState turn on/off a specific server.
//turnOn=true to turn on, turnOn=false to turn off
func (c *Client) setServerPowerState(id string, powerState bool) error {
	isOn, err := c.IsServerOn(id)
	if err != nil {
		return err
	}
	if isOn == powerState {
		return nil
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "power"),
		method: http.MethodPatch,
		body: ServerPowerUpdateRequest{
			Power: powerState,
		},
	}
	err = r.execute(*c, nil)
	if err != nil {
		return err
	}
	return c.WaitForServerPowerStatus(id, powerState)
}

//StartServer starts a server
func (c *Client) StartServer(id string) error {
	return c.setServerPowerState(id, true)
}

//StopServer stops a server
func (c *Client) StopServer(id string) error {
	return c.setServerPowerState(id, false)
}

//ShutdownServer shutdowns a specific server
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
		uri:    path.Join(apiServerBase, id, "shutdown"),
		method: http.MethodPatch,
		body:   map[string]string{},
	}

	err = r.execute(*c, nil)
	if err != nil {
		if requestError, ok := err.(RequestError); ok {
			if requestError.StatusCode == 500 {
				c.cfg.logger.Debugf("Graceful shutdown for server %s has failed. power-off will be used", id)
				return c.StopServer(id)
			}
		}
		return err
	}

	//If we get an error, which includes a timeout, power off the server instead
	err = c.WaitForServerPowerStatus(id, false)
	if err != nil {
		c.cfg.logger.Debugf("Graceful shutdown for server %s has failed. power-off will be used", id)
		return c.StopServer(id)
	}
	return nil
}

//GetServersByLocation gets a list of servers by location
func (c *Client) GetServersByLocation(id string) ([]Server, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "servers"),
		method: http.MethodGet,
	}
	var response ServerList
	var servers []Server
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{Properties: properties})
	}
	return servers, err
}

//GetDeletedServers gets a list of deleted servers
func (c *Client) GetDeletedServers() ([]Server, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "servers"),
		method: http.MethodGet,
	}
	var response DeletedServerList
	var servers []Server
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		servers = append(servers, Server{Properties: properties})
	}
	return servers, err
}
