package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetObjectStorageAccessKeyList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "access_keys")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareObjectStorageAccessKeyListHTTPGet())
	})

	res, err := client.GetObjectStorageAccessKeyList()
	assert.Nil(t, err, "GetObjectStorageAccessKeyList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockObjectStorageAccessKey()), fmt.Sprintf("%v", res))
}

func TestClient_GetObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "access_keys", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareObjectStorageAccessKeyHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		res, err := client.GetObjectStorageAccessKey(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetObjectStorageAccessKey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockObjectStorageAccessKey()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_CreateObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	var isFailed bool
	uri := path.Join(apiObjectStorageBase, "access_keys")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodPost, request.Method)
		if isFailed {
			writer.WriteHeader(400)
		} else {
			fmt.Fprint(writer, prepareObjectStorageAccessKeyHTTPCreateResponse())
		}
	})
	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})
	for _, test := range commonSuccessFailTestCases {
		isFailed = test.isFailed
		res, err := client.CreateObjectStorageAccessKey()
		if isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteObjectStorageAccessKey returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockObjectStorageAccessKeyCreateResponse()), fmt.Sprintf("%v", res))
		}
	}
}

func TestClient_DeleteObjectStorageAccessKey(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "access_keys", dummyUUID)
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodDelete, request.Method)
		fmt.Fprint(writer, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.DeleteObjectStorageAccessKey(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteObjectStorageAccessKey returned an error %v", err)
		}
	}
}

func TestClient_GetObjectStorageBucketList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiObjectStorageBase, "buckets")
	mux.HandleFunc(uri, func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, http.MethodGet, request.Method)
		fmt.Fprintf(writer, prepareObjectStorageBucketListHTTPGet())
	})

	res, err := client.GetObjectStorageBucketList()
	assert.Nil(t, err, "GetObjectStorageBucketList returned an error %v", err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockObjectStorageBucket()), fmt.Sprintf("%v", res))
}

func getMockObjectStorageAccessKey() ObjectStorageAccessKey {
	mock := ObjectStorageAccessKey{Properties: ObjectStorageAccessKeyProperties{
		SecretKey: "dummy-secret-key",
		AccessKey: "dummy-access-key",
		User:      "any",
	}}
	return mock
}

func prepareObjectStorageAccessKeyListHTTPGet() string {
	accessKey := getMockObjectStorageAccessKey()
	res, _ := json.Marshal(accessKey.Properties)
	return fmt.Sprintf(`{"access_keys": [%s]}`, string(res))
}

func prepareObjectStorageAccessKeyHTTPGet() string {
	accessKey := getMockObjectStorageAccessKey()
	res, _ := json.Marshal(accessKey)
	return string(res)
}

func getMockObjectStorageAccessKeyCreateResponse() ObjectStorageAccessKeyCreateResponse {
	mock := ObjectStorageAccessKeyCreateResponse{
		AccessKey: struct {
			SecretKey string `json:"secret_key"`
			AccessKey string `json:"access_key"`
		}{
			SecretKey: "dummy-secret-key",
			AccessKey: "dummy-access-key",
		},
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareObjectStorageAccessKeyHTTPCreateResponse() string {
	createRes := getMockObjectStorageAccessKeyCreateResponse()
	res, _ := json.Marshal(createRes)
	return string(res)
}

func getMockObjectStorageBucket() ObjectStorageBucket {
	mock := ObjectStorageBucket{Properties: ObjectStorageBucketProperties{
		Name: "test",
		Usage: struct {
			SizeKb     int `json:"size_kb"`
			NumObjects int `json:"num_objects"`
		}{
			SizeKb:     1000000,
			NumObjects: 10,
		},
	}}
	return mock
}

func prepareObjectStorageBucketListHTTPGet() string {
	bucket := getMockObjectStorageBucket()
	res, _ := json.Marshal(bucket.Properties)
	return fmt.Sprintf(`{"buckets": [%s]}`, string(res))
}
