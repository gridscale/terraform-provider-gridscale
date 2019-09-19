package gsclient

import (
	"errors"
	"net/http"
	"path"
)

//StorageSnapshotList is JSON structure of a list of storage snapshots
type StorageSnapshotList struct {
	List map[string]StorageSnapshotProperties `json:"snapshots"`
}

//DeletedStorageSnapshotList is JSON structure of a list of deleted storage snapshots
type DeletedStorageSnapshotList struct {
	List map[string]StorageSnapshotProperties `json:"deleted_snapshots"`
}

//StorageSnapshot is JSON structure of a single storage snapshot
type StorageSnapshot struct {
	Properties StorageSnapshotProperties `json:"snapshot"`
}

//StorageSnapshotProperties JSON struct of properties of a storage snapshot
type StorageSnapshotProperties struct {
	Labels           []string `json:"labels"`
	ObjectUUID       string   `json:"object_uuid"`
	Name             string   `json:"name"`
	Status           string   `json:"status"`
	LocationCountry  string   `json:"location_country"`
	UsageInMinutes   int      `json:"usage_in_minutes"`
	LocationUUID     string   `json:"location_uuid"`
	ChangeTime       string   `json:"change_time"`
	LicenseProductNo int      `json:"license_product_no"`
	CurrentPrice     float64  `json:"current_price"`
	CreateTime       string   `json:"create_time"`
	Capacity         int      `json:"capacity"`
	LocationName     string   `json:"location_name"`
	LocationIata     string   `json:"location_iata"`
	ParentUUID       string   `json:"parent_uuid"`
}

//StorageSnapshotCreateRequest JSON struct of a request for creating a storage snapshot
type StorageSnapshotCreateRequest struct {
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//StorageSnapshotCreateResponse JSON struct of a response for creating a storage snapshot
type StorageSnapshotCreateResponse struct {
	RequestUUID string `json:"request_uuid"`
	ObjectUUID  string `json:"object_uuid"`
}

//StorageSnapshotUpdateRequest JSON struct of a request for updating a storage snapshot
type StorageSnapshotUpdateRequest struct {
	Name   string   `json:"name,omitempty"`
	Labels []string `json:"labels,omitempty"`
}

//StorageRollbackRequest JSON struct of a request for rolling back
type StorageRollbackRequest struct {
	Rollback bool `json:"rollback,omitempty"`
}

//StorageSnapshotExportToS3Request JSON struct of a request for exporting a storage snapshot to S3
type StorageSnapshotExportToS3Request struct {
	S3auth struct {
		Host      string `json:"host"`
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
	} `json:"s3auth"`
	S3data struct {
		Host     string `json:"host"`
		Bucket   string `json:"bucket"`
		Filename string `json:"filename"`
		Private  bool   `json:"private"`
	} `json:"s3data"`
}

//GetStorageSnapshotList gets a list of storage snapshots
func (c *Client) GetStorageSnapshotList(id string) ([]StorageSnapshot, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshots"),
		method: http.MethodGet,
	}
	var response StorageSnapshotList
	var snapshots []StorageSnapshot
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		snapshots = append(snapshots, StorageSnapshot{Properties: properties})
	}
	return snapshots, err
}

//GetStorageSnapshot gets a specific storage's snapshot based on given storage id and snapshot id.
func (c *Client) GetStorageSnapshot(storageID, snapshotID string) (StorageSnapshot, error) {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return StorageSnapshot{}, errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID),
		method: http.MethodGet,
	}
	var response StorageSnapshot
	err := r.execute(*c, &response)
	return response, err
}

//CreateStorageSnapshot creates a new storage's snapshot
func (c *Client) CreateStorageSnapshot(id string, body StorageSnapshotCreateRequest) (StorageSnapshotCreateResponse, error) {
	if !isValidUUID(id) {
		return StorageSnapshotCreateResponse{}, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, id, "snapshots"),
		method: http.MethodPost,
		body:   body,
	}
	var response StorageSnapshotCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return StorageSnapshotCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//UpdateStorageSnapshot updates a specific storage's snapshot
func (c *Client) UpdateStorageSnapshot(storageID, snapshotID string, body StorageSnapshotUpdateRequest) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//DeleteStorageSnapshot deletes a specific storage's snapshot
func (c *Client) DeleteStorageSnapshot(storageID, snapshotID string) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//RollbackStorage rollbacks a storage
func (c *Client) RollbackStorage(storageID, snapshotID string, body StorageRollbackRequest) error {
	if !isValidUUID(storageID) || !isValidUUID(snapshotID) {
		return errors.New("'storageID' or 'snapshotID' is invalid")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID, "rollback"),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//ExportStorageSnapshotToS3 export a storage's snapshot to S3
func (c *Client) ExportStorageSnapshotToS3(storageID, snapshotID string, body StorageSnapshotExportToS3Request) error {
	if storageID == "" || snapshotID == "" {
		return errors.New("'storageID' and 'snapshotID' are required")
	}
	r := Request{
		uri:    path.Join(apiStorageBase, storageID, "snapshots", snapshotID, "export_to_s3"),
		method: http.MethodPatch,
		body:   body,
	}
	return r.execute(*c, nil)
}

//GetSnapshotsByLocation gets a list of storage snapshots by location
func (c *Client) GetSnapshotsByLocation(id string) ([]StorageSnapshot, error) {
	if !isValidUUID(id) {
		return nil, errors.New("'id' is invalid")
	}
	r := Request{
		uri:    path.Join(apiLocationBase, id, "snapshots"),
		method: http.MethodGet,
	}
	var response StorageSnapshotList
	var snapshots []StorageSnapshot
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		snapshots = append(snapshots, StorageSnapshot{Properties: properties})
	}
	return snapshots, err
}

//GetDeletedSnapshots gets a list of deleted storage snapshots
func (c *Client) GetDeletedSnapshots() ([]StorageSnapshot, error) {
	r := Request{
		uri:    path.Join(apiDeletedBase, "snapshots"),
		method: http.MethodGet,
	}
	var response DeletedStorageSnapshotList
	var snapshots []StorageSnapshot
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		snapshots = append(snapshots, StorageSnapshot{Properties: properties})
	}
	return snapshots, err
}
