package gsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	uri				string
	method			string
	body			interface{}
}

type CreateResponse struct {
	ObjectUuid  string `json:"object_uuid"`
	RequestUuid string `json:"request_uuid"`
	ServerUuid	string	`json:"server_uuid"`
}

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *Request) execute(c Client, output interface{}) (error) {
	url := c.cfg.APIUrl + r.uri

	//Convert the body of the request to json
	jsonBody := new(bytes.Buffer)
	if r.body != nil {
		err := json.NewEncoder(jsonBody).Encode(r.body)
		if err != nil{
			return err
		}
	}

	//Add authentication headers and content type
	request, err := http.NewRequest(r.method, url, jsonBody)
	if err != nil{
		return err
	}
	request.Header.Add("X-Auth-UserId", c.cfg.UserUUID)
	request.Header.Add("X-Auth-Token", c.cfg.APIToken)
	request.Header.Add("Content-Type", "application/json")

	log.Printf("[DEBUG] Request body: %v", request.Body)

	//execute the request
	result, err := c.cfg.HTTPClient.Do(request)
	if err != nil{
		return err
	}

	iostream, err := ioutil.ReadAll(result.Body)
	if err != nil{
		return err
	}
	json.Unmarshal(iostream, output) //Edit the given struct
	response := string(iostream)

	log.Printf("[DEBUG] Response body: %v", response)

	if result.StatusCode >= 300 {
		return fmt.Errorf("[Error] statuscode %v returned", result.StatusCode)
	}

	return nil
}