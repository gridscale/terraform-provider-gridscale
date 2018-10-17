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

func (r *Request) execute(c Client, s interface{}) (  error) {
	url := c.cfg.APIUrl + r.uri

	jsonBody, _ := json.Marshal(r.body)

	request, err := http.NewRequest(r.method, url, bytes.NewBuffer(jsonBody))
	request.Header.Add("X-Auth-UserId", c.cfg.UserUUID)
	request.Header.Add("X-Auth-Token", c.cfg.APIToken)
	request.Header.Add("Content-Type", "application/json")

	log.Printf("[DEBUG] Request body: %v", request.Body)

	result, err := c.cfg.HTTPClient.Do(request)

	iostream, err := ioutil.ReadAll(result.Body)

	json.Unmarshal(iostream, s)

	body := string(iostream)

	log.Printf("[DEBUG] Response body: %v", body)

	return err
}