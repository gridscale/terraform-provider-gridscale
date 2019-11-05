package gsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//request gridscale's custom request struct
type request struct {
	uri                 string
	method              string
	body                interface{}
	skipCheckingRequest bool
}

//CreateResponse common struct of a response for creation
type CreateResponse struct {
	//UUID of the object being created
	ObjectUUID string `json:"object_uuid"`

	//UUID of the request
	RequestUUID string `json:"request_uuid"`
}

//RequestStatus status of a request
type RequestStatus map[string]RequestStatusProperties

//RequestStatusProperties JSON struct of properties of a request's status
type RequestStatusProperties struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	CreateTime GSTime `json:"create_time"`
}

//RequestError error of a request
type RequestError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusCode  int
	RequestUUID string
}

//Error just returns error as string
func (r RequestError) Error() string {
	message := r.Description
	if message == "" {
		message = "no error message received from server"
	}
	errorMessageFormat := "Status code: %v. Error: %s. Request UUID: %s. "
	if r.StatusCode >= 500 {
		errorMessageFormat += "Please report this error along with the request UUID."
	}
	return fmt.Sprintf(errorMessageFormat, r.StatusCode, message, r.RequestUUID)
}

const requestUUIDHeaderParam = "X-Request-Id"

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *request) execute(ctx context.Context, c Client, output interface{}) error {
	url := c.cfg.apiURL + r.uri
	logger := c.Logger()
	httpClient := c.HttpClient()

	logger.Debugf("%v request sent to URL: %v", r.method, url)

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
	request = request.WithContext(ctx)
	request.Header.Set("User-Agent", c.UserAgent())
	request.Header.Add("X-Auth-UserID", c.UserUUID())
	request.Header.Add("X-Auth-Token", c.APIToken())
	request.Header.Add("Content-Type", bodyType)
	logger.Debugf("Request body: %v", request.Body)

	//Init request UUID variable
	var requestUUID string
	//Init empty response body
	var responseBodyBytes []byte
	//
	err = retryWithLimitedNumOfRetries(func() (bool, error) {
		//execute the request
		result, err := httpClient.Do(request)
		if err != nil {
			logger.Errorf("Error while executing the request: %v", err)
			return false, err
		}
		statusCode := result.StatusCode
		requestUUID = result.Header.Get(requestUUIDHeaderParam)
		responseBodyBytes, err = ioutil.ReadAll(result.Body)
		if err != nil {
			logger.Errorf("Error while reading the response's body: %v", err)
			return false, err
		}

		logger.Debugf("Status code: %v. Request UUID: %v.", statusCode, requestUUID)

		if result.StatusCode >= 300 {
			var errorMessage RequestError //error messages have a different structure, so they are read with a different struct
			errorMessage.StatusCode = statusCode
			errorMessage.RequestUUID = requestUUID
			json.Unmarshal(responseBodyBytes, &errorMessage)
			//if internal server error OR object is in status that does not allow the request, retry
			if result.StatusCode >= 500 || result.StatusCode == 424 {
				return true, errorMessage
			}
			logger.Errorf(
				"Error message: %v. Title: %v. Code: %v. Request UUID: %v.",
				errorMessage.Description,
				errorMessage.Title,
				errorMessage.StatusCode,
				errorMessage.RequestUUID,
			)
			return false, errorMessage
		}
		logger.Debugf("Response body: %v", string(responseBodyBytes))

		return false, nil
	}, c.MaxNumberOfRetries(), c.DelayInterval())
	//if retry fails
	if err != nil {
		return err
	}

	//if output is set
	if output != nil {
		//Unmarshal body bytes to the given struct
		err = json.Unmarshal(responseBodyBytes, output)
		if err != nil {
			logger.Errorf("Error while marshaling JSON: %v", err)
			return err
		}
	}

	//If the client is synchronous, and the request does not skip
	//checking a request, wait until the request completes
	if c.Synchronous() && !r.skipCheckingRequest {
		return c.waitForRequestCompleted(ctx, requestUUID)
	}
	return nil
}
