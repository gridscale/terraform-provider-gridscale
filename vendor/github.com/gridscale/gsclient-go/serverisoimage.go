package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ServerIsoImageRelationList JSON struct of a list of relations between a server and ISO-Images
type ServerIsoImageRelationList struct {
	List []ServerIsoImageRelationProperties `json:"isoimage_relations"`
}

//ServerIsoImageRelation JSON struct of a single relation between a server and an ISO-Image
type ServerIsoImageRelation struct {
	Properties ServerIsoImageRelationProperties `json:"isoimage_relation"`
}

//ServerIsoImageRelationProperties JSON struct of properties of a relation between a server and an ISO-Image
type ServerIsoImageRelationProperties struct {
	ObjectUUID string `json:"object_uuid"`
	ObjectName string `json:"object_name"`
	Private    bool   `json:"private"`
	CreateTime string `json:"create_time"`
	Bootdevice bool   `json:"bootdevice"`
}

//ServerIsoImageRelationCreateRequest JSON struct of a request for creating a relation between a server and an ISO-Image
type ServerIsoImageRelationCreateRequest struct {
	ObjectUUID string `json:"object_uuid"`
}

//ServerIsoImageRelationUpdateRequest JSON struct of a request for updating a relation between a server and an ISO-Image
type ServerIsoImageRelationUpdateRequest struct {
	BootDevice bool   `json:"bootdevice"`
	Name       string `json:"name"`
}

//GetServerIsoImageList gets a list of a specific server's ISO images
func (c *Client) GetServerIsoImageList(id string) ([]ServerIsoImageRelationProperties, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "isoimages"),
		method: http.MethodGet,
	}
	var response ServerIsoImageRelationList
	err := r.execute(*c, &response)
	return response.List, err
}

//GetServerIsoImage gets an ISO image of a specific server
func (c *Client) GetServerIsoImage(serverID, isoImageID string) (ServerIsoImageRelationProperties, error) {
	if !isValidUUID(serverID) || !isValidUUID(isoImageID) {
		return ServerIsoImageRelationProperties{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "isoimages", isoImageID),
		method: http.MethodGet,
	}
	var response ServerIsoImageRelation
	err := r.execute(*c, &response)
	return response.Properties, err
}

//UpdateServerIsoImage updates a link between a storage and an ISO image
func (c *Client) UpdateServerIsoImage(serverID, isoImageID string, body ServerIsoImageRelationUpdateRequest) error {
	if !isValidUUID(serverID) || !isValidUUID(isoImageID) {
		return errors.New("'serverID' or 'isoImageID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "isoimages", isoImageID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//CreateServerIsoImage creates a link between a server and an ISO image
func (c *Client) CreateServerIsoImage(id string, body ServerIsoImageRelationCreateRequest) error {
	if !isValidUUID(id) || !isValidUUID(body.ObjectUUID) {
		return errors.New("'serverID' or 'isoImageID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, id, "isoimages"),
		method: http.MethodPost,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteServerIsoImage deletes a link between an ISO image and a server
func (c *Client) DeleteServerIsoImage(serverID, isoImageID string) error {
	if !isValidUUID(serverID) || !isValidUUID(isoImageID) {
		return errors.New("'serverID' or 'isoImageID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiServerBase, serverID, "isoimages", isoImageID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//LinkIsoImage attaches an ISO image to a server
func (c *Client) LinkIsoImage(serverID string, isoimageID string) error {
	body := ServerIsoImageRelationCreateRequest{
		ObjectUUID: isoimageID,
	}
	return c.CreateServerIsoImage(serverID, body)
}

//UnlinkIsoImage removes the link between an ISO image and a server
func (c *Client) UnlinkIsoImage(serverID string, isoimageID string) error {
	return c.DeleteServerIsoImage(serverID, isoimageID)
}
