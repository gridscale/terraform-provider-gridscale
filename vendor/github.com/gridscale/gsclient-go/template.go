package gsclient

import (
	"fmt"
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

func (c *Client) GetTemplate(id string) (*Template, error) {
	r := Request{
		uri:    apiTemplateBase + "/" + id,
		method: "GET",
	}

	response := new(Template)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) GetTemplateList() ([]Template, error) {
	r := Request{
		uri:    apiTemplateBase,
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
