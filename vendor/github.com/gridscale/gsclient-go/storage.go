package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//StorageList JSON struct of a list of storages
type StorageList struct {
	List map[string]StorageProperties `json:"storages"`
}

//DeletedStorageList JSON struct of a list of storages
type DeletedStorageList struct {
	List map[string]StorageProperties `json:"deleted_storages"`
}

//Storage JSON struct of a single storage
type Storage struct {
	Properties StorageProperties `json:"storage"`
}

//StorageProperties JSON struct of properties of a storage
type StorageProperties struct {
	ChangeTime       string                    `json:"change_time"`
	LocationIata     string                    `json:"location_iata"`
	Status           string                    `json:"status"`
	LicenseProductNo int                       `json:"license_product_no"`
	LocationCountry  string                    `json:"location_country"`
	UsageInMinutes   int                       `json:"usage_in_minutes"`
	LastUsedTemplate string                    `json:"last_used_template"`
	CurrentPrice     float64                   `json:"current_price"`
	Capacity         int                       `json:"capacity"`
	LocationUUID     string                    `json:"location_uuid"`
	StorageType      string                    `json:"storage_type"`
	ParentUUID       string                    `json:"parent_uuid"`
	Name             string                    `json:"name"`
	LocationName     string                    `json:"location_name"`
	ObjectUUID       string                    `json:"object_uuid"`
	Snapshots        []StorageSnapshotRelation `json:"snapshots"`
	Relations        StorageRelations          `json:"relations"`
	Labels           []string                  `json:"labels"`
	CreateTime       string                    `json:"create_time"`
}

//StorageRelations JSON struct of a list of a storage's relations
type StorageRelations struct {
	Servers           []StorageServerRelation              `json:"servers"`
	SnapshotSchedules []StorageAndSnapshotScheduleRelation `json:"snapshot_schedules"`
}

//StorageServerRelation JSON struct of a relation between a storage and a server
type StorageServerRelation struct {
	Bootdevice bool   `json:"bootdevice"`
	Target     int    `json:"target"`
	Controller int    `json:"controller"`
	Bus        int    `json:"bus"`
	ObjectUUID string `json:"object_uuid"`
	Lun        int    `json:"lun"`
	CreateTime string `json:"create_time"`
	ObjectName string `json:"object_name"`
}

//StorageSnapshotRelation JSON struct of a relation between a storage and a snapshot
type StorageSnapshotRelation struct {
	LastUsedTemplate      string `json:"last_used_template"`
	ObjectUUID            string `json:"object_uuid"`
	StorageUUID           string `json:"storage_uuid"`
	SchedulesSnapshotName string `json:"schedules_snapshot_name"`
	SchedulesSnapshotUUID string `json:"schedules_snapshot_uuid"`
	ObjectCapacity        int    `json:"object_capacity"`
	CreateTime            string `json:"create_time"`
	ObjectName            string `json:"object_name"`
}

//StorageAndSnapshotScheduleRelation JSON struct of a relation between a storage and a snapshot schedule
type StorageAndSnapshotScheduleRelation struct {
	RunInterval   int    `json:"run_interval"`
	KeepSnapshots int    `json:"keep_snapshots"`
	ObjectName    string `json:"object_name"`
	NextRuntime   string `json:"next_runtime"`
	ObjectUUID    int    `json:"object_uuid"`
	Name          string `json:"name"`
	CreateTime    string `json:"create_time"`
}

//StorageTemplate JSON struct of a storage template
type StorageTemplate struct {
	Sshkeys      []string `json:"sshkeys,omitempty"`
	TemplateUUID string   `json:"template_uuid"`
	Password     string   `json:"password,omitempty"`
	PasswordType string   `json:"password_type,omitempty"`
	Hostname     string   `json:"hostname,omitempty"`
}

//StorageCreateRequest JSON struct of a request for creating a storage
type StorageCreateRequest struct {
	Capacity     int              `json:"capacity"`
	LocationUUID string           `json:"location_uuid"`
	Name         string           `json:"name"`
	StorageType  string           `json:"storage_type,omitempty"`
	Template     *StorageTemplate `json:"template,omitempty"`
	Labels       []string         `json:"labels,omitempty"`
}

//StorageUpdateRequest JSON struct of a request for updating a storage
type StorageUpdateRequest struct {
	Name     string   `json:"name,omitempty"`
	Labels   []string `json:"labels,omitempty"`
	Capacity int      `json:"capacity,omitempty"`
}

//GetStorage get a storage
func (c *Client) GetStorage(id string) (Storage, error) {
	if !isValidUUID(id) {
		return Storage{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodGet,
	}
	var response Storage
	err := r.execute(*c, &response)
	return response, err
}

//GetStorageList gets a list of available storages
func (c *Client) GetStorageList() ([]Storage, error) {
	r := Request{
		uri:    apiStorageBase,
		method: http.MethodGet,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{
			Properties: properties,
		})
	}
	return storages, err
}

//CreateStorage create a storage
func (c *Client) CreateStorage(body StorageCreateRequest) (CreateResponse, error) {
	r := Request{
		uri:    apiStorageBase,
		method: http.MethodPost,
		body:   body,
	}
	var response CreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return CreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//DeleteStorage delete a storage
func (c *Client) DeleteStorage(id string) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//UpdateStorage update a storage
func (c *Client) UpdateStorage(id string, body StorageUpdateRequest) error {
	if !isValidUUID(id) {
		return errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetStorageEventList get list of a storage's events
func (c *Client) GetStorageEventList(id string) ([]Event, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "events"),
		method: http.MethodGet,
	}
	var response EventList
	var storageEvents []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		storageEvents = append(storageEvents, Event{Properties: properties})
	}
	return storageEvents, err
}

//GetStoragesByLocation gets a list of storages by location
func (c *Client) GetStoragesByLocation(id string) ([]Storage, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "storages"),
		method: http.MethodGet,
	}
	var response StorageList
	var storages []Storage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{Properties: properties})
	}
	return storages, err
}

//GetDeletedStorages gets a list of deleted storages
func (c *Client) GetDeletedStorages() ([]Storage, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "storages"),
		method: http.MethodGet,
	}
	var response DeletedStorageList
	var storages []Storage
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		storages = append(storages, Storage{Properties: properties})
	}
	return storages, err
}
