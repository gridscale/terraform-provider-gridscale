package gsclient

import (
	"errors"
	"net/http"
	"path"
	"strings"
)

//ObjectStorageAccessKeyList is JSON structure of a list of Object Storage Access Keys
type ObjectStorageAccessKeyList struct {
	List []ObjectStorageAccessKeyProperties `json:"access_keys"`
}

//ObjectStorageAccessKey is JSON structure of a single Object Storage Access Key
type ObjectStorageAccessKey struct {
	Properties ObjectStorageAccessKeyProperties `json:"access_key"`
}

//ObjectStorageAccessKeyProperties is JSON struct of properties of an object storage access key
type ObjectStorageAccessKeyProperties struct {
	SecretKey string `json:"secret_key"`
	AccessKey string `json:"access_key"`
	User      string `json:"user"`
}

//ObjectStorageAccessKeyCreateResponse is JSON struct of a response for creating an object storage access key
type ObjectStorageAccessKeyCreateResponse struct {
	AccessKey struct {
		SecretKey string `json:"secret_key"`
		AccessKey string `json:"access_key"`
	} `json:"access_key"`
	RequestUUID string `json:"request_uuid"`
}

//ObjectStorageBucketList is JSON struct of a list of buckets
type ObjectStorageBucketList struct {
	List []ObjectStorageBucketProperties `json:"buckets"`
}

//ObjectStorageBucket is JSON struct of a single bucket
type ObjectStorageBucket struct {
	Properties ObjectStorageBucketProperties `json:"bucket"`
}

//ObjectStorageBucketProperties is JSON struct of properties of a bucket
type ObjectStorageBucketProperties struct {
	Name  string `json:"name"`
	Usage struct {
		SizeKb     int `json:"size_kb"`
		NumObjects int `json:"num_objects"`
	} `json:"usage"`
}

//GetObjectStorageAccessKeyList gets a list of available object storage access keys
func (c *Client) GetObjectStorageAccessKeyList() ([]ObjectStorageAccessKey, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys"),
		method: http.MethodGet,
	}
	var response ObjectStorageAccessKeyList
	var accessKeys []ObjectStorageAccessKey
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		accessKeys = append(accessKeys, ObjectStorageAccessKey{Properties: properties})
	}
	return accessKeys, err
}

//GetObjectStorageAccessKey gets a specific object storage access key based on given id
func (c *Client) GetObjectStorageAccessKey(id string) (ObjectStorageAccessKey, error) {
	if strings.TrimSpace(id) == "" {
		return ObjectStorageAccessKey{}, errors.New("'id' is required")
	}
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys", id),
		method: http.MethodGet,
	}
	var response ObjectStorageAccessKey
	err := r.execute(*c, &response)
	return response, err
}

//CreateObjectStorageAccessKey creates an object storage access key
func (c *Client) CreateObjectStorageAccessKey() (ObjectStorageAccessKeyCreateResponse, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys"),
		method: http.MethodPost,
	}
	var response ObjectStorageAccessKeyCreateResponse
	err := r.execute(*c, &response)
	if err != nil {
		return ObjectStorageAccessKeyCreateResponse{}, err
	}
	err = c.WaitForRequestCompletion(response.RequestUUID)
	return response, err
}

//DeleteObjectStorageAccessKey deletes a specific object storage access key based on given id
func (c *Client) DeleteObjectStorageAccessKey(id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("'id' is required")
	}
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "access_keys", id),
		method: http.MethodDelete,
	}
	return r.execute(*c, nil)
}

//GetObjectStorageBucketList gets a list of object storage buckets
func (c *Client) GetObjectStorageBucketList() ([]ObjectStorageBucket, error) {
	r := Request{
		uri:    path.Join(apiObjectStorageBase, "buckets"),
		method: http.MethodGet,
	}
	var response ObjectStorageBucketList
	var buckets []ObjectStorageBucket
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		buckets = append(buckets, ObjectStorageBucket{Properties: properties})
	}
	return buckets, err
}
