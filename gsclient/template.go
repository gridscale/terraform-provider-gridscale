package gsclient

import (
	"fmt"
	"log"
)

type Templates struct {
	List map[string]TemplateProperties `json:"templates"`
}

type Template struct {
	Properties TemplateProperties `json:"template"`
}

type TemplateProperties struct {
	Name       string `json:"name"`
	ObjectUuid string `json:"object_uuid"`
	Status     string `json:"status"`
	CreateTime string `json:"create_time"`
	ChangeTime string `json:"change_time"`
	Template   string `json:"sshkey"`
}

type TemplateCreateRequest struct {
	Name         string   `json:"name"`
	Labels       []string `json:"labels"`
	SnapshotUuid string   `json:"name"`
}

func (c *Client) GetTemplate(id string) (*Template, error) {
	r := Request{
		uri:    "/objects/templates/" + id,
		method: "GET",
	}

	response := new(Template)
	err := r.execute(*c, &response)

	log.Printf("Received sshkey: %v", response)

	return response, err
}

func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    "/objects/templates/",
		method: "GET",
	}

	response := new(Templates)
	err := r.execute(*c, &response)

	list := []Template{}
	for _, properties := range response.List {
		template := Template{
			Properties: properties,
		}
		list = append(list, template)
	}

	return list, err
}

func (c *Client) CreateTemplate(body TemplateCreateRequest) (*CreateResponse, error) {
	r := Request{
		uri:    "/objects/templates",
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

func (c *Client) DeleteTemplate(id string) error {
	r := Request{
		uri:    "/objects/template/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateTemplate(id string, body map[string]interface{}) error {
	r := Request{
		uri:    "/objects/sshkeys/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}

func (c *Client) GetTemplateByName(name string) (*Template, error) {
	templates, err := c.GetTemplateList()
	if err != nil {
		return nil, err
	}
	for _, template := range templates {
		if template.Properties.Name == name {
			return c.GetTemplate(template.Properties.ObjectUuid)
		}
	}

	return nil, fmt.Errorf("Template %v not found", name)
}
