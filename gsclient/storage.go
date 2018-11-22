package gsclient

type Storages struct {
	List map[string]StorageProperties `json:"storages"`
}

type Storage struct {
	Properties StorageProperties `json:"storage"`
}

type StorageProperties struct {
	ChangeTime       string            `json:"change_time"`
	LocationIata     string            `json:"location_iata"`
	Status           string            `json:"status"`
	LicenseProductNo int               `json:"license_product_no"`
	LocationCountry  string            `json:"location_country"`
	UsageInMinutes   int               `json:"usage_in_minutes"`
	LastUsedTemplate string            `json:"last_used_template"`
	CurrentPrice     float64           `json:"current_price"`
	Capacity         int               `json:"capacity"`
	LocationUuid     string            `json:"location_uuid"`
	StorageType      string            `json:"storage_type"`
	ParentUuid       string            `json:"parent_uuid"`
	Name             string            `json:"name"`
	LocationName     string            `json:"location_name"`
	ObjectUuid       string            `json:"object_uuid"`
	Snapshots        []StorageSnapshot `json:"snapshots"`
	Relations        StorageRelations  `json:"relations"`
	Labels           []string          `json:"labels"`
	CreateTime       string            `json:"create_time"`
}

type StorageRelations struct {
	Servers           []StorageServer            `json:"servers"`
	SnapshotSchedules []StorageSnapshotSchedules `json:"snapshot_schedules"`
}

type StorageServer struct {
	Bootdevice bool   `json:"bootdevice"`
	Target     int    `json:"target"`
	Controller int    `json:"controller"`
	Bus        int    `json:"bus"`
	ObjectUuid string `json:"object_uuid"`
	Lun        int    `json:"lun"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
}

type StorageSnapshot struct {
	LastUsedTemplate      string `json:"last_used_template"`
	ObjectUuid            string `json:"object_uuid"`
	SchedulesSnapshotName string `json:"schedules_snapshot_name"`
	SchedulesSnapshotUuid string `json:"schedules_snapshot_uuid"`
	ObjectCapacity        int    `json:"object_capacity"`
	CreateTime            string `json:"create_time"`
	ObjectName            string `json:"object_name"`
}

type StorageSnapshotSchedules struct {
	RunInterval   int    `json:"run_interval"`
	KeepSnapshots int    `json:"keep_snapshots"`
	ObjectName    string `json:"object_name"`
	NextRuntime   string `json:"next_runtime"`
	ObjectUuid    int    `json:"object_uuid"`
	Name          string `json:"name"`
	CreateTime    string `json:"create_time"`
}
type StorageTemplate struct {
	Sshkeys      []string `json:"sshkeys,omitempty"`
	TemplateUuid string   `json:"template_uuid"`
	Password     string   `json:"password,omitempty"`
	PasswordType string   `json:"password_type,omitempty"`
	Hostname     string   `json:"hostname,omitempty"`
}

type StorageCreateRequest struct {
	Capacity     int             `json:"capacity"`
	LocationUuid string          `json:"location_uuid"`
	Name         string          `json:"name"`
	StorageType  string          `json:"storage_type,omitempty"`
	Template     StorageTemplate `json:"template,omitempty"`
	Labels       []interface{}   `json:"labels,omitempty"`
}

type StorageUpdateRequest struct {
	Name     string        `json:"name,omitempty"`
	Labels   []interface{} `json:"labels,omitempty"`
	Capacity int           `json:"capacity"`
}

func (c *Client) GetStorage(id string) (*Storage, error) {
	r := Request{
		uri:    apiStorageBase + "/" + id,
		method: "GET",
	}

	response := new(Storage)
	err := r.execute(*c, &response)

	return response, err
}

func (c *Client) GetStorageList() ([]Storage, error) {
	r := Request{
		uri:    apiStorageBase,
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

func (c *Client) CreateStorage(body StorageCreateRequest) (*CreateResponse, error) {
	r := Request{
		uri:    apiStorageBase,
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
		uri:    apiStorageBase + "/" + id,
		method: "DELETE",
	}

	return r.execute(*c, nil)
}

func (c *Client) UpdateStorage(id string, body StorageUpdateRequest) error {
	r := Request{
		uri:    apiStorageBase + "/" + id,
		method: "PATCH",
		body:   body,
	}

	return r.execute(*c, nil)
}
