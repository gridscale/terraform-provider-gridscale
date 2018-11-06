package gsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Request struct {
	uri    string
	method string
	body   interface{}
}

type CreateResponse struct {
	ObjectUuid  string `json:"object_uuid"`
	RequestUuid string `json:"request_uuid"`
	ServerUuid  string `json:"server_uuid"`
}

type RequestStatus map[string]RequestStatusProperties

type RequestStatusProperties struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime string `json:"create_time"`
}

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *Request) execute(c Client, output interface{}) error {
	url := c.cfg.APIUrl + r.uri

	//Convert the body of the request to json
	jsonBody := new(bytes.Buffer)
	if r.body != nil {
		err := json.NewEncoder(jsonBody).Encode(r.body)
		if err != nil {
			return err
		}
	}

	//Add authentication headers and content type
	request, err := http.NewRequest(r.method, url, jsonBody)
	if err != nil {
		return err
	}
	request.Header.Add("X-Auth-UserId", c.cfg.UserUUID)
	request.Header.Add("X-Auth-Token", c.cfg.APIToken)
	request.Header.Add("Content-Type", "application/json")

	log.Printf("[DEBUG] Request body: %v", request.Body)

	//execute the request
	result, err := c.cfg.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	iostream, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(iostream, output) //Edit the given struct
	response := string(iostream)

	log.Printf("[DEBUG] Response body: %v", response)

	if result.StatusCode >= 300 {
		return fmt.Errorf("[Error] statuscode %v returned: %v", result.StatusCode, result.Body)
	}

	return nil
}

//This function allows use to wait for a request to complete. Timeouts are currently hardcoded
func (c *Client) WaitForRequestCompletion(cr CreateResponse) error {
	r := Request{
		uri:    "/requests/" + cr.RequestUuid,
		method: "GET",
	}

	timer := time.After(30 * time.Second)

	for {
		select {
		case <-timer:
			return fmt.Errorf("Timeout reached when waiting for request %v to complete", cr.RequestUuid)
		default:
			time.Sleep(500 * time.Millisecond) //delay the request, so we don't do too many requests to the server
			response := new(RequestStatus)
			r.execute(*c, &response)
			output := *response //Without this cast reading indexes doesn't work
			if output[cr.RequestUuid].Status == "done" {
				log.Print("Done with creating")
				return nil
			}
		}
	}
}

func (c *Client) WaitForServerPowerStatus(id string, status bool) error {
	timer := time.After(30 * time.Second)

	for {
		select {
		case <-timer:
			return fmt.Errorf("Timeout reached when trying to shut down system with id %v", id)
		default:
			time.Sleep(500 * time.Millisecond) //delay the request, so we don't do too many requests to the server
			server, err := c.GetServer(id)
			if err != nil {
				return err
			}
			if server.Properties.Power == status {
				log.Print("The power status of the server with id %v has changed to %s", id, status)
				return nil
			}
		}
	}
}
