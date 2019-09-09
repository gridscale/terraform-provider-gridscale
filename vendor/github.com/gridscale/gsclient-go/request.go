package gsclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

//Request gridscale's custom request struct
type Request struct {
	uri    string
	method string
	body   interface{}
}

//CreateResponse common struct of a response for creation
type CreateResponse struct {
	ObjectUUID  string `json:"object_uuid"`
	RequestUUID string `json:"request_uuid"`
}

//RequestStatus status of a request
type RequestStatus map[string]RequestStatusProperties

//RequestStatusProperties JSON struct of properties of a request's status
type RequestStatusProperties struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime string `json:"create_time"`
}

//RequestError error of a request
type RequestError struct {
	StatusMessage string `json:"status"`
	ErrorMessage  string `json:"message"`
	StatusCode    int
}

//Error just returns error as string
func (r RequestError) Error() string {
	message := r.ErrorMessage
	if message == "" {
		message = "no error message received from server"
	}
	return fmt.Sprintf("statuscode %v returned: %s", r.StatusCode, message)
}

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *Request) execute(c Client, output interface{}) error {
	url := c.cfg.APIUrl + r.uri
	c.cfg.logger.Debugf("%v request sent to URL: %v", r.method, url)

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
	request.Header.Add("X-Auth-UserID", c.cfg.UserUUID)
	request.Header.Add("X-Auth-Token", c.cfg.APIToken)
	request.Header.Add("Content-Type", "application/json")
	c.cfg.logger.Debugf("Request body: %v", request.Body)

	//execute the request
	result, err := c.cfg.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	iostream, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return err
	}

	c.cfg.logger.Debugf("Status code returned: %v", result.StatusCode)

	if result.StatusCode >= 300 {
		var errorMessage RequestError //error messages have a different structure, so they are read with a different struct
		errorMessage.StatusCode = result.StatusCode
		json.Unmarshal(iostream, &errorMessage)
		c.cfg.logger.Errorf("Error message: %v. Status: %v. Code: %v.", errorMessage.ErrorMessage, errorMessage.StatusMessage, errorMessage.StatusCode)
		return errorMessage
	}
	json.Unmarshal(iostream, output) //Edit the given struct
	c.cfg.logger.Debugf("Response body: %v", string(iostream))
	return nil
}

//WaitForRequestCompletion allows to wait for a request to complete. Timeouts are currently hardcoded
func (c *Client) WaitForRequestCompletion(id string) error {
	r := Request{
		uri:    path.Join("/requests/", id),
		method: "GET",
	}
	timer := time.After(time.Minute)

	for {
		select {
		case <-timer:
			c.cfg.logger.Errorf("Timeout reached when waiting for request %v to complete", id)
			return fmt.Errorf("Timeout reached when waiting for request %v to complete", id)
		default:
			time.Sleep(500 * time.Millisecond) //delay the request, so we don't do too many requests to the server
			var response RequestStatus
			r.execute(*c, &response)
			if response[id].Status == "done" {
				c.cfg.logger.Info("Done with creating")
				return nil
			}
		}
	}
}

//WaitForServerPowerStatus  allows to wait for a server changing its power status. Timeouts are currently hardcoded
func (c *Client) WaitForServerPowerStatus(id string, status bool) error {
	timer := time.After(2 * time.Minute)
	for {
		select {
		case <-timer:
			c.cfg.logger.Errorf("Timeout reached when trying to shut down system with id %v", id)
			return fmt.Errorf("Timeout reached when trying to shut down system with id %v", id)
		default:
			time.Sleep(500 * time.Millisecond) //delay the request, so we don't do too many requests to the server
			server, err := c.GetServer(id)
			if err != nil {
				return err
			}
			if server.Properties.Power == status {
				c.cfg.logger.Infof("The power status of the server with id %v has changed to %t", id, status)
				return nil
			}
		}
	}
}
