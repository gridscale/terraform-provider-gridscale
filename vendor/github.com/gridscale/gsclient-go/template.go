package gsclient

import (
	"errors"
	"fmt"
	"net/http"
	"path"
)

//TemplateList JSON struct of a list of templates
type TemplateList struct {
	List map[string]TemplateProperties `json:"templates"`
}

//DeletedTemplateList JSON struct of a list of deleted templates
type DeletedTemplateList struct {
	List map[string]TemplateProperties `json:"deleted_templates"`
}

//Template JSON struct of a single template
type Template struct {
	Properties TemplateProperties `json:"template"`
}

//TemplateProperties JSOn struct of properties of a template
type TemplateProperties struct {
	Status           string   `json:"status"`
	Ostype           string   `json:"ostype"`
	LocationUUID     string   `json:"location_uuid"`
	Version          string   `json:"version"`
	LocationIata     string   `json:"location_iata"`
	ChangeTime       string   `json:"change_time"`
	Private          bool     `json:"private"`
	ObjectUUID       string   `json:"object_uuid"`
	LicenseProductNo int      `json:"license_product_no"`
	CreateTime       string   `json:"create_time"`
	UsageInMinutes   int      `json:"usage_in_minutes"`
	Capacity         int      `json:"capacity"`
	LocationName     string   `json:"location_name"`
	Distro           string   `json:"distro"`
	Description      string   `json:"description"`
	CurrentPrice     float64  `json:"current_price"`
	LocationCountry  string   `json:"location_country"`
	Name             string   `json:"name"`
	Labels           []string `json:"labels"`
}

//TemplateCreateRequest JSON struct of a request for creating a template
type TemplateCreateRequest struct {
	Name         string   `json:"name"`
	SnapshotUUID string   `json:"snapshot_uuid"`
	Labels       []string `json:"labels,omitempty"`
}

//TemplateUpdateRequest JSON struct of a request for updating a template
type TemplateUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//GetTemplate gets a template
func (c *Client) GetTemplate(id string) (Template, error) {
	if !isValidUUID(id) {
		return Template{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodGet,
	}
	var response Template
	err := r.execute(*c, &response)
	return response, err
}

//GetTemplateList gets a list of templates
func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{
			Properties: properties,
		})
	}
	return templates, err
}

//GetTemplateByName gets a template by its name
func (c *Client) GetTemplateByName(name string) (Template, error) {
	if name == "" {
		return Template{}, errors.New("'name' is required")
	}
	templates, err := c.GetTemplateList()
	if err != nil {
		return Template{}, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return Template{Properties: template.Properties}, nil
		}
	}
	return Template{}, fmt.Errorf("Template %v not found", name)
}

//CreateTemplate creates a template
func (c *Client) CreateTemplate(body TemplateCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiTemplateBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(*c, &response)
	return response, err
}

//UpdateTemplate updates a template
func (c *Client) UpdateTemplate(id string, body TemplateUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteTemplate deletes a template
func (c *Client) DeleteTemplate(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetTemplateEventList gets a list of a template's events
func (c *Client) GetTemplateEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiTemplateBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var templateEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templateEvents = append(templateEvents, Event{Properties: properties})
	}
	return templateEvents, err
}

//GetTemplatesByLocation gets a list of templates by location
func (c *Client) GetTemplatesByLocation(id string) ([]Template, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "templates"),
		method: http.MethodGet,
	}
	var response TemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{Properties: properties})
	}
	return templates, err
}

//GetDeletedTemplates gets a list of deleted templates
func (c *Client) GetDeletedTemplates() ([]Template, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "templates"),
		method: http.MethodGet,
	}
	var response DeletedTemplateList
	var templates []Template
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		templates = append(templates, Template{Properties: properties})
	}
	return templates, err
}
