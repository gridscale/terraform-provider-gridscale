package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetStorageList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageListHTTPGet())
	})
	response, err := client.GetStorageList()
	assert.Nil(t, err, "GetStorageList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage()), fmt.Sprintf("%v", response))
}

func TestClient_GetStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetStorage(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorage returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockStorage()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_CreateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiStorageBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareFirewallCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateStorage(StorageCreateRequest{
		Capacity:     10,
		LocationUUID: dummyUUID,
		Name:         "test",
		StorageType:  "storage",
		Template: &StorageTemplate{
			TemplateUUID: dummyUUID,
			Password:     "pass",
			PasswordType: "crypt",
			Hostname:     "example.com",
		},
		Labels: []string{"label"},
	})
	assert.Nil(t, err, "CreateStorage returned an error %v", err)
	assert.Equal(t, fmt.Sprintf("%v", getMockStorageCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.UpdateStorage(test.testUUID, StorageUpdateRequest{
			Name:     "test",
			Labels:   []string{"label"},
			Capacity: 20,
		})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "UpdateStorage returned an error %v", err)
		}
	}
}

func TestClient_DeleteStorage(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.DeleteStorage(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteStorage returned an error %v", err)
		}
	}
}

func TestClient_GetStorageEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiStorageBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetStorageEventList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStorageEventList returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_GetStoragesByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "storages")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareStorageListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetStoragesByLocation(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetStoragesByLocation returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_GetDeletedStorages(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiDeletedBase, "storages")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareDeletedStorageListHTTPGet())
	})
	response, err := client.GetDeletedStorages()
	assert.Nil(t, err, "GetDeletedStorages returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockStorage()), fmt.Sprintf("%v", response))
}

func getMockStorage() Storage {
	mock := Storage{Properties: StorageProperties{
		ChangeTime:       dummyTime,
		LocationIata:     "iata",
		Status:           "active",
		LicenseProductNo: 11111,
		LocationCountry:  "Germany",
		UsageInMinutes:   10,
		LastUsedTemplate: dummyUUID,
		CurrentPrice:     9.1,
		Capacity:         10,
		LocationUUID:     dummyUUID,
		StorageType:      "storage",
		ParentUUID:       dummyUUID,
		Name:             "test",
		LocationName:     "Cologne",
		ObjectUUID:       dummyUUID,
		Snapshots: []StorageSnapshotRelation{
			{
				LastUsedTemplate:      dummyUUID,
				ObjectUUID:            dummyUUID,
				StorageUUID:           dummyUUID,
				SchedulesSnapshotName: "test",
				SchedulesSnapshotUUID: dummyUUID,
				ObjectCapacity:        10,
				CreateTime:            dummyTime,
				ObjectName:            "test",
			},
		},
		Relations:  StorageRelations{},
		Labels:     []string{"label"},
		CreateTime: dummyTime,
	}}
	return mock
}

func getMockStorageCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareStorageListHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"storages": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareStorageHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage)
	return string(res)
}

func prepareStorageCreateResponse() string {
	response := getMockStorageCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareDeletedStorageListHTTPGet() string {
	storage := getMockStorage()
	res, _ := json.Marshal(storage.Properties)
	return fmt.Sprintf(`{"deleted_storages": {"%s": %s}}`, dummyUUID, string(res))
}
