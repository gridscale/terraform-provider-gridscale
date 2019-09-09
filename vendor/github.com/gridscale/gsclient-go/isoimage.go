package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//ISOImageList is JSON struct of a list of ISO images
type ISOImageList struct {
	List map[string]ISOImageProperties `json:"isoimages"`
}

//DeletedISOImageList is JSON struct of a list of deleted SO images
type DeletedISOImageList struct {
	List map[string]ISOImageProperties `json:"deleted_isoimages"`
}

//ISOImage is JSON struct of a list an ISO image
type ISOImage struct {
	Properties ISOImageProperties `json:"isoimage"`
}

//ISOImageProperties is JSON struct of properties of an ISO image
type ISOImageProperties struct {
	ObjectUUID      string           `json:"object_uuid"`
	Relations       ISOImageRelation `json:"relations"`
	Description     string           `json:"description"`
	LocationName    string           `json:"location_name"`
	SourceURL       string           `json:"source_url"`
	Labels          []string         `json:"labels"`
	LocationIata    string           `json:"location_iata"`
	LocationUUID    string           `json:"location_uuid"`
	Status          string           `json:"status"`
	CreateTime      string           `json:"create_time"`
	Name            string           `json:"name"`
	Version         string           `json:"version"`
	LocationCountry string           `json:"location_country"`
	UsageInMinutes  int              `json:"usage_in_minutes"`
	Private         bool             `json:"private"`
	ChangeTime      string           `json:"change_time"`
	Capacity        int              `json:"capacity"`
	CurrentPrice    float64          `json:"current_price"`
}

//ISOImageRelation is JSON struct of a list of an ISO-Image's relations
type ISOImageRelation struct {
	Servers []ServerinISOImage `json:"servers"`
}

//ServerinISOImage is JSON struct of a relation between an ISO-Image and a Server
type ServerinISOImage struct {
	Bootdevice bool   `json:"bootdevice"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
	ObjectUUID string `json:"object_uuid"`
}

//ISOImageCreateRequest is JSON struct of a request for creating an ISO-Image
type ISOImageCreateRequest struct {
	Name         string   `json:"name"`
	SourceURL    string   `json:"source_url"`
	Labels       []string `json:"labels,omitempty"`
	LocationUUID string   `json:"location_uuid"`
}

//ISOImageCreateResponse is JSON struct of a response for creating an ISO-Image
type ISOImageCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
}

//ISOImageUpdateRequest is JSON struct of a request for updating an ISO-Image
type ISOImageUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//GetISOImageList returns a list of available ISO images
func (c *Client) GetISOImageList() ([]ISOImage, error) {
	r := Request{
		uri:    path.Join(apiISOBase),
		method: http.MethodGet,
	}
	var response ISOImageList
	var isoImages []ISOImage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		isoImages = append(isoImages, ISOImage{Properties: properties})
	}
	return isoImages, err
}

//GetISOImage returns a specific ISO image based on given id
func (c *Client) GetISOImage(id string) (ISOImage, error) {
	if !isValidUUID(id) {
		return ISOImage{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodGet,
	}
	var response ISOImage
	err := r.execute(*c, &response)
	return response, err
}

//CreateISOImage creates an ISO image
func (c *Client) CreateISOImage(body ISOImageCreateRequest) (ISOImageCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiISOBase),
		method: http.MethodPost,
		body:   body,
	}
	var response ISOImageCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return ISOImageCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//UpdateISOImage updates a specific ISO Image
func (c *Client) UpdateISOImage(id string, body ISOImageUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteISOImage deletes a specific ISO image
func (c *Client) DeleteISOImage(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetISOImageEventList returns a list of events of an ISO image
func (c *Client) GetISOImageEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiISOBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var isoImageEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		isoImageEvents = append(isoImageEvents, Event{Properties: properties})
	}
	return isoImageEvents, err
}

//GetISOImagesByLocation gets a list of ISO images by location
func (c *Client) GetISOImagesByLocation(id string) ([]ISOImage, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "isoimages"),
		method: http.MethodGet,
	}
	var response ISOImageList
	var isoImages []ISOImage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		isoImages = append(isoImages, ISOImage{Properties: properties})
	}
	return isoImages, err
}

//GetDeletedISOImages gets a list of deleted ISO images
func (c *Client) GetDeletedISOImages() ([]ISOImage, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "isoimages"),
		method: http.MethodGet,
	}
	var response DeletedISOImageList
	var isoImages []ISOImage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		isoImages = append(isoImages, ISOImage{Properties: properties})
	}
	return isoImages, err
}
