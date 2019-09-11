package gsclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path"
	"testing"
)

func TestClient_GetTemplateList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	response, err := client.GetTemplateList()
	assert.Nil(t, err, "GetTemplateList returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate()), fmt.Sprintf("%v", response))
}

func TestClient_GetTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplate(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplate returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplate()), fmt.Sprintf("%v", response))
		}

	}
}

func TestClient_GetTemplateByName(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	testCases := []uuidTestCase{
		{
			testUUID: "test",
			isFailed: false,
		},
		{
			testUUID: "",
			isFailed: true,
		},
	}
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	for _, test := range testCases {
		response, err := client.GetTemplateByName(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplateByName returned an error %v", err)
			assert.Equal(t, fmt.Sprintf("%v", getMockTemplate()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_CreateTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := apiTemplateBase
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		fmt.Fprintf(w, prepareTemplateCreateResponse())
	})

	httpResponse := fmt.Sprintf(`{"%s": {"status":"done"}}`, dummyRequestUUID)
	mux.HandleFunc("/requests/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, httpResponse)
	})

	res, err := client.CreateTemplate(TemplateCreateRequest{
		Name:         "test",
		SnapshotUUID: dummyUUID,
		Labels:       []string{"label"},
	})
	assert.Nil(t, err, "CreateTemplate returned an error %v", err)
	assert.Equal(t, fmt.Sprintf("%v", getMockTemplateCreateResponse()), fmt.Sprintf("%v", res))
}

func TestClient_UpdateTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		fmt.Fprintf(w, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.UpdateTemplate(test.testUUID, TemplateUpdateRequest{
			Name:   "test",
			Labels: []string{"labels"},
		})
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "UpdateTemplate returned an error %v", err)
		}
	}
}

func TestClient_DeleteTemplate(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID)
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		fmt.Fprintf(w, "")
	})
	for _, test := range uuidCommonTestCases {
		err := client.DeleteTemplate(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "DeleteTemplate returned an error %v", err)
		}
	}
}

func TestClient_GetTemplateEventList(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiTemplateBase, dummyUUID, "events")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareEventListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplateEventList(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplateEventList returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockEvent()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_GetTemplatesByLocation(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiLocationBase, dummyUUID, "templates")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareTemplateListHTTPGet())
	})
	for _, test := range uuidCommonTestCases {
		response, err := client.GetTemplatesByLocation(test.testUUID)
		if test.isFailed {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err, "GetTemplatesByLocation returned an error %v", err)
			assert.Equal(t, 1, len(response))
			assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate()), fmt.Sprintf("%v", response))
		}
	}
}

func TestClient_GetDeletedTemplates(t *testing.T) {
	server, client, mux := setupTestClient()
	defer server.Close()
	uri := path.Join(apiDeletedBase, "templates")
	mux.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		fmt.Fprint(w, prepareDeletedTemplateListHTTPGet())
	})
	response, err := client.GetDeletedTemplates()
	assert.Nil(t, err, "GetDeletedTemplates returned an error %v", err)
	assert.Equal(t, 1, len(response))
	assert.Equal(t, fmt.Sprintf("[%v]", getMockTemplate()), fmt.Sprintf("%v", response))
}

func getMockTemplate() Template {
	mock := Template{Properties: TemplateProperties{
		Status:           "active",
		Ostype:           "type",
		LocationUUID:     dummyUUID,
		Version:          "1.0",
		LocationIata:     "iata",
		ChangeTime:       dummyTime,
		Private:          true,
		ObjectUUID:       dummyUUID,
		LicenseProductNo: 11111,
		CreateTime:       dummyTime,
		UsageInMinutes:   1000,
		Capacity:         10,
		LocationName:     "Cologne",
		Distro:           "Centos7",
		Description:      "description",
		CurrentPrice:     0,
		LocationCountry:  "Germnany",
		Name:             "test",
		Labels:           []string{"label"},
	}}
	return mock
}

func getMockTemplateCreateResponse() CreateResponse {
	mock := CreateResponse{
		ObjectUUID:  dummyUUID,
		RequestUUID: dummyRequestUUID,
	}
	return mock
}

func prepareTemplateListHTTPGet() string {
	template := getMockTemplate()
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"templates": {"%s": %s}}`, dummyUUID, string(res))
}

func prepareTemplateHTTPGet() string {
	template := getMockTemplate()
	res, _ := json.Marshal(template)
	return string(res)
}

func prepareTemplateCreateResponse() string {
	response := getMockTemplateCreateResponse()
	res, _ := json.Marshal(response)
	return string(res)
}

func prepareDeletedTemplateListHTTPGet() string {
	template := getMockTemplate()
	res, _ := json.Marshal(template.Properties)
	return fmt.Sprintf(`{"deleted_templates": {"%s": %s}}`, dummyUUID, string(res))
}
