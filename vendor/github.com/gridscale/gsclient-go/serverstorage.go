package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ServerStorageRelationList JSON struct of a list of relations between a server and storages
type ServerStorageRelationList struct {
	List []ServerStorageRelationProperties `json:"storage_relations"`
}

//ServerStorageRelationSingle JSON struct of a single relation between a server and a storage
type ServerStorageRelationSingle struct {
	Properties ServerStorageRelationProperties `json:"storage_relation"`
}

//ServerStorageRelationProperties JSON struct of properties of a relation between a server and a storage
type ServerStorageRelationProperties struct {
	ObjectUUID       string `json:"object_uuid"`
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
	ServerUUID       string `json:"server_uuid"`
}

//ServerStorageRelationCreateRequest JSON struct of a request for creating a relation between a server and a storage
type ServerStorageRelationCreateRequest struct {
	ObjectUUID string `json:"object_uuid"`
	BootDevice bool   `json:"bootdevice,omitempty"`
}

//ServerStorageRelationUpdateRequest JSON struct of a request for updating a relation between a server and a storage
type ServerStorageRelationUpdateRequest struct {
	Ordering   int      `json:"ordering,omitempty"`
	BootDevice bool     `json:"bootdevice,omitempty"`
	L3security []string `json:"l3security,omitempty"`
}

//GetServerStorageList gets a list of a specific server's storages
func (c *Client) GetServerStorageList(id string) ([]ServerStorageRelationProperties, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "storages"),
		method: http.MethodGet,
	}
	var response ServerStorageRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerStorage gets a storage of a specific server
func (c *Client) GetServerStorage(serverID, storageID string) (ServerStorageRelationProperties, error) {
	if !isValidUUID(serverID) || !isValidUUID(storageID) {
		return ServerStorageRelationProperties{}, errors.New("'serverID' or 'storageID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "storages", storageID),
		method: http.MethodGet,
	}
	var response ServerStorageRelationSingle
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerStorage updates a link between a storage and a server
func (c *Client) UpdateServerStorage(serverID, storageID string, body ServerStorageRelationUpdateRequest) error {
	if !isValidUUID(serverID) || !isValidUUID(storageID) {
		return errors.New("'serverID' or 'storageID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "storages", storageID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerStorage create a link between a server and a storage
func (c *Client) CreateServerStorage(id string, body ServerStorageRelationCreateRequest) error {
	if !isValidUUID(id) || !isValidUUID(body.ObjectUUID) {
		return errors.New("'server_id' or 'storage_id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "storages"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerStorage delete a link between a storage and a server
func (c *Client) DeleteServerStorage(serverID, storageID string) error {
	if !isValidUUID(serverID) || !isValidUUID(storageID) {
		return errors.New("'serverID' or 'storageID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "storages", storageID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkStorage attaches a storage to a server
func (c *Client) LinkStorage(serverID string, storageID string, bootdevice bool) error {
	body := ServerStorageRelationCreateRequest{
		ObjectUUID: storageID,
		BootDevice: bootdevice,
	}
	return c.CreateServerStorage(serverID, body)
}

//UnlinkStorage remove a storage from a server
func (c *Client) UnlinkStorage(serverID string, storageID string) error {
	return c.DeleteServerStorage(serverID, storageID)
}
