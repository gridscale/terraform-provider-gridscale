package gsclient

type Storages struct {
	List map[string]StorageProperties `json:"storages"`
}

type Storage struct {
	Properties StorageProperties `json:"storage"`
}

type StorageProperties struct {
	ObjectUuid      string  `json:"object_uuid"`
	Name            string  `json:"name"`
	Capacity        string  `json:"capacity"`
	LocationUuid    string  `json:"location_uuid"`
	Status          string  `json:"status"`
	CreateTime      string  `json:"create_time"`
	ChangeTime      string  `json:"change_time"`
	CurrentPrice    float64 `json:"current_price"`
	LocationName    string  `json:"location_name"`
	LocationCountry string  `json:"location_country"`
	LocationIata    string  `json:"location_iata"`
}

type StorageCreateRequest struct {
	Capacity     int             `json:"capacity"`
	LocationUuid string          `json:"location_uuid"`
	Name         string          `json:"name"`
	StorageType  string          `json:"storage_type"`
	Template     StorageTemplate `json:"template,omitempty"`
}

type StorageTemplate struct {
	Sshkeys      []string `json:"sshkeys"`
	TemplateUuid string   `json:"template_uuid"`
}

func (c *Client) GetStorage(id string) (*Storage, error) {
	r := Request{
		uri:    "/objects/storages/" + id,
		method: "GET",
	}

	response := new(Storage)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) GetStorageList() ([]Storage, error) {
	r := Request{
		uri:    "/objects/storages/",
		method: "GET",
	}

	response := new(Storages)
	err := r.execute(*c, &response)

	list := []Storage{}
	for _, properties := range response.List {
		storage := Storage{
			Properties: properties,
		}
		list = append(list, storage)
	}

	return list, err
}

func (c *Client) CreateStorage(body map[string]interface{}) (*CreateResponse, error) {
	r := Request{
		uri:    "/objects/storages",
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

func (c *Client) DeleteStorage(id string) error {
	r := Request{
		uri:    "/objects/storages/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateStorage(id string, body map[string]interface{}) error {
	r := Request{
		uri:    "/objects/storages/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}
