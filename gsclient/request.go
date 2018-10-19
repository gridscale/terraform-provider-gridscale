package gsclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	uri				string
	method			string
	body			map[string]interface{}
}

type CreateResponse struct {
	ObjectUuid  string `json:"object_uuid"`
	RequestUuid string `json:"request_uuid"`
}

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *Request) execute(c Client, s interface{}) (  error) {
	url := c.cfg.APIUrl + r.uri

	//Convert the body of the request to json
	jsonBody, _ := json.Marshal(r.body)

	//Add authentication headers and content type
	request, err := http.NewRequest(r.method, url, bytes.NewBuffer(jsonBody))
	request.Header.Add("X-Auth-UserId", c.cfg.UserUUID)
	request.Header.Add("X-Auth-Token", c.cfg.APIToken)
	request.Header.Add("Content-Type", "application/json")

	log.Printf("[DEBUG] Request body: %v", request.Body)

	//execute the request
	result, err := c.cfg.HTTPClient.Do(request)

	iostream, err := ioutil.ReadAll(result.Body)
	json.Unmarshal(iostream, s) //Edit the given struct
	body := string(iostream)

	log.Printf("[DEBUG] Response body: %v", body)

	return err
}