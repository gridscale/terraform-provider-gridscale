package gsclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

//gsRequest gridscale's custom gsRequest struct
type gsRequest struct {
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

const (
	requestUUIDHeaderParam           = "X-Request-Id"
	requestRateLimitResetHeaderParam = "ratelimit-reset"
)

//This function takes the client and a struct and then adds the result to the given struct if possible
func (r *gsRequest) execute(ctx context.Context, c Client, output interface{}) error {

	//Prepare http request (including HTTP headers preparation, etc.)
	httpReq, err := r.prepareHTTPRequest(ctx, c.cfg)
	logger.Debugf("Request body: %v", httpReq.Body)
	logger.Debugf("Request headers: %v", httpReq.Header)

	//Execute the request (including retrying when needed)
	requestUUID, responseBodyBytes, err := r.retryHTTPRequest(ctx, c.HttpClient(), httpReq, c.MaxNumberOfRetries(), c.DelayInterval())
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

//prepareHTTPRequest prepares a http request
func (r *gsRequest) prepareHTTPRequest(ctx context.Context, cfg *Config) (*http.Request, error) {
	url := cfg.apiURL + r.uri
	logger.Debugf("Preparing %v request sent to URL: %v", r.method, url)

	//Convert the body of the request to json
	jsonBody := new(bytes.Buffer)
	if r.body != nil {
		err := json.NewEncoder(jsonBody).Encode(r.body)
		if err != nil {
			return nil, err
		}
	}

	//Add authentication headers and content type
	request, err := http.NewRequest(r.method, url, jsonBody)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	request.Header.Set("User-Agent", cfg.userAgent)
	request.Header.Set("X-Auth-UserID", cfg.userUUID)
	request.Header.Set("X-Auth-Token", cfg.apiToken)
	request.Header.Set("Content-Type", bodyType)

	//Set headers based on a given list of custom headers
	//Use Header.Set() instead of Header.Add() because we want to
	//override the headers' values if they are already set.
	for k, v := range cfg.httpHeaders {
		request.Header.Set(k, v)
	}

	return request, nil
}

//retryHTTPRequest runs & retries a HTTP request
//returns UUID (string), response body ([]byte), error
func (r *gsRequest) retryHTTPRequest(
	ctx context.Context,
	httpClient *http.Client,
	httpReq *http.Request,
	maxNoOfRetries int,
	delayInterval time.Duration,
) (string, []byte, error) {
	//Init request UUID variable
	var requestUUID string
	//Init empty response body
	var responseBodyBytes []byte
	//
	err := retryWithLimitedNumOfRetries(func() (bool, error) {
		//execute the request
		resp, err := httpClient.Do(httpReq)
		if err != nil {
			//If the error is caused by expired context, return context error and no need to retry
			if ctx.Err() != nil {
				return false, ctx.Err()
			}

			if err, ok := err.(net.Error); ok {
				// exclude retry request with none GET method (write operations) in case of a request timeout or a context error
				if err.Timeout() && r.method != http.MethodGet {
					return false, err
				}
				logger.Debugf("Retrying request due to network error %v", err)
				return true, err
			}
			logger.Errorf("Error while executing the request: %v", err)
			//stop retrying (false) and return error
			return false, err
		}
		//Close body to prevent resource leak
		defer resp.Body.Close()

		statusCode := resp.StatusCode
		requestUUID = resp.Header.Get(requestUUIDHeaderParam)
		responseBodyBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("Error while reading the response's body: %v", err)
			//stop retrying (false) and return error
			return false, err
		}

		logger.Debugf("Status code: %v. Request UUID: %v. Headers: %v", statusCode, requestUUID, resp.Header)

		//If the status code is an error code
		if statusCode >= 300 {
			var errorMessage RequestError //error messages have a different structure, so they are read with a different struct
			errorMessage.StatusCode = statusCode
			errorMessage.RequestUUID = requestUUID
			json.Unmarshal(responseBodyBytes, &errorMessage)

			//If the error is retryable
			if isErrorHTTPCodeRetryable(statusCode) {
				//If status code is 429, that means we reach the rate limit
				if statusCode == 429 {
					//Get the time that the rate limit will be reset
					rateLimitResetTimestamp := resp.Header.Get(requestRateLimitResetHeaderParam)
					//Get the delay time
					delayMs, err := getDelayTimeInMsFromTimestampStr(rateLimitResetTimestamp)
					if err != nil {
						return false, err
					}
					//Delay the retry until the rate limit is reset
					logger.Debugf("Delay request for %d ms: %v method sent to url %v with body %v", delayMs, r.method, httpReq.URL.RequestURI(), r.body)
					time.Sleep(time.Duration(delayMs) * time.Millisecond)
				}
				logger.Debugf("Retrying request: %v method sent to url %v with body %v", r.method, httpReq.URL.RequestURI(), r.body)
				return true, errorMessage
			}

			//If the error is not retryable
			logger.Errorf(
				"Error message: %v. Title: %v. Code: %v. Request UUID: %v.",
				errorMessage.Description,
				errorMessage.Title,
				errorMessage.StatusCode,
				errorMessage.RequestUUID,
			)
			//stop retrying (false) and return custom error
			return false, errorMessage

		}
		logger.Debugf("Response body: %v", string(responseBodyBytes))
		//stop retrying (false) as no more errors
		return false, nil
	}, maxNoOfRetries, delayInterval)
	//No need to return when the context is already expired.
	select {
	case <-ctx.Done():
		return "", nil, ctx.Err()
	default:
	}
	return requestUUID, responseBodyBytes, err
}

func isErrorHTTPCodeRetryable(statusCode int) bool {
	//if internal server error (>=500)
	//OR object is in status that does not allow the request (424)
	//OR we reach the rate limit (429), retry
	if statusCode >= 500 || statusCode == 424 || statusCode == 429 {
		return true
	}
	//stop retrying (false) and return custom error
	return false
}

//getDelayTimeInMsFromTimestampStr takes a unix timestamp (ms) string
//and returns the amount of ms to delay the next HTTP request retry
//return error if the input string is not a valid unix timestamp(ms)
func getDelayTimeInMsFromTimestampStr(timestamp string) (int64, error) {
	if timestamp == "" {
		return 0, errors.New("timestamp is empty")
	}
	//convert timestamp from string to int
	timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return 0, err
	}
	currentTimestampMs := time.Now().UnixNano() / 1000000
	return timestampInt - currentTimestampMs, nil
}
