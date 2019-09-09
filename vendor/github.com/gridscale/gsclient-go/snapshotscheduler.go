package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//StorageSnapshotScheduleList JSON of a list of storage snapshot schedule
type StorageSnapshotScheduleList struct {
	List map[string]StorageSnapshotScheduleProperties `json:"snapshot_schedules"`
}

//StorageSnapshotSchedule JSON struct of a single storage snapshot schedule
type StorageSnapshotSchedule struct {
	Properties StorageSnapshotScheduleProperties `json:"snapshot_schedule"`
}

//StorageSnapshotScheduleProperties JSON struct of properties of a single storage snapshot schedule
type StorageSnapshotScheduleProperties struct {
	ChangeTime    string                           `json:"change_time"`
	CreateTime    string                           `json:"create_time"`
	KeepSnapshots int                              `json:"keep_snapshots"`
	Labels        []string                         `json:"labels"`
	Name          string                           `json:"name"`
	NextRuntime   string                           `json:"next_runtime"`
	ObjectUUID    string                           `json:"object_uuid"`
	Relations     StorageSnapshotScheduleRelations `json:"relations"`
	RunInterval   int                              `json:"run_interval"`
	Status        string                           `json:"status"`
	StorageUUID   string                           `json:"storage_uuid"`
}

//StorageSnapshotScheduleRelations JSON struct of a list of relations of a storage snapshot schedule
type StorageSnapshotScheduleRelations struct {
	Snapshots []StorageSnapshotScheduleRelation `json:"snapshots"`
}

//StorageSnapshotScheduleRelation JSON struct of a relation of a storage snapshot schedule
type StorageSnapshotScheduleRelation struct {
	CreateTime string `json:"create_time"`
	Name       string `json:"name"`
	ObjectUUID string `json:"object_uuid"`
}

//StorageSnapshotScheduleCreateRequest JSON struct of a request for creating a storage snapshot schedule
type StorageSnapshotScheduleCreateRequest struct {
	Name          string   `json:"name"`
	Labels        []string `json:"labels,omitempty"`
	RunInterval   int      `json:"run_interval"`
	KeepSnapshots int      `json:"keep_snapshots"`
	NextRuntime   string   `json:"next_runtime,omitempty"`
}

//StorageSnapshotScheduleCreateResponse JSON struct of a response for creating a storage snapshot schedule
type StorageSnapshotScheduleCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
}

//StorageSnapshotScheduleUpdateRequest JSON struct of a request for updating a storage snapshot schedule
type StorageSnapshotScheduleUpdateRequest struct {
	Name          string   `json:"name,omitempty"`
	Labels        []string `json:"labels,omitempty"`
	RunInterval   int      `json:"run_interval,omitempty"`
	KeepSnapshots int      `json:"keep_snapshots,omitempty"`
	NextRuntime   string   `json:"next_runtime,omitempty"`
}

//GetStorageSnapshotScheduleList gets a list of available storage snapshot schedules based on a given storage's id
func (c *Client) GetStorageSnapshotScheduleList(id string) ([]StorageSnapshotSchedule, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshot_schedules"),
		method: http.MethodGet,
	}
	var response StorageSnapshotScheduleList
	var schedules []StorageSnapshotSchedule
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		schedules = append(schedules, StorageSnapshotSchedule{Properties: properties})
	}
	return schedules, err
}

//GetStorageSnapshotSchedule gets a specific storage snapshot scheduler based on a given storage's id and scheduler's id
func (c *Client) GetStorageSnapshotSchedule(storageID, scheduleID string) (StorageSnapshotSchedule, error) {
	if !isValidUUID(storageID) || !isValidUUID(scheduleID) {
		return StorageSnapshotSchedule{}, errors.New("'storageID' or 'scheduleID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshot_schedules", scheduleID),
		method: http.MethodGet,
	}
	var response StorageSnapshotSchedule
	err := r.execute(*c, &response)
	return response, err
}

//CreateStorageSnapshotSchedule create a storage's snapshot scheduler
func (c *Client) CreateStorageSnapshotSchedule(id string, body StorageSnapshotScheduleCreateRequest) (
	StorageSnapshotScheduleCreateResponse, error) {
	if !isValidUUID(id) {
		return StorageSnapshotScheduleCreateResponse{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshot_schedules"),
		method: http.MethodPost,
		body:   body,
	}
	var response StorageSnapshotScheduleCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return StorageSnapshotScheduleCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//UpdateStorageSnapshotSchedule updates specific Storage's snapshot scheduler based on a given storage's id and scheduler's id
func (c *Client) UpdateStorageSnapshotSchedule(storageID, scheduleID string,
	body StorageSnapshotScheduleUpdateRequest) error {
	if !isValidUUID(storageID) || !isValidUUID(scheduleID) {
		return errors.New("'storageID' or 'scheduleID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshot_schedules", scheduleID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteStorageSnapshotSchedule deletes specific Storage's snapshot scheduler based on a given storage's id and scheduler's id
func (c *Client) DeleteStorageSnapshotSchedule(storageID, scheduleID string) error {
	if !isValidUUID(storageID) || !isValidUUID(scheduleID) {
		return errors.New("'storageID' or 'scheduleID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshot_schedules", scheduleID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}
