package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//LocationList JSON struct of a list of locations
type LocationList struct {
	List map[string]LocationProperties `json:"locations"`
}

//Location JSON struct of a single location
type Location struct {
	Properties LocationProperties `json:"location"`
}

//LocationProperties JSON struct of properties of a location
type LocationProperties struct {
	Iata       string   `json:"iata"`
	Status     string   `json:"status"`
	Labels     []string `json:"labels"`
	Name       string   `json:"name"`
	ObjectUUID string   `json:"object_uuid"`
	Country    string   `json:"country"`
}

//GetLocationList gets a list of available locations
func (c *Client) GetLocationList() ([]Location, error) {
	r := Request{
		uri:    apiLocationBase,
		method: http.MethodGet,
	}
	var response LocationList
	var locations []Location
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		locations = append(locations, Location{Properties: properties})
	}
	return locations, err
}

//GetLocation gets a specific location
func (c *Client) GetLocation(id string) (Location, error) {
	if !isValidUUID(id) {
		return Location{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id),
		method: http.MethodGet,
	}
	var location Location
	err := r.execute(*c, &location)
	return location, err
}
