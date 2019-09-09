package gsclient

import "net/http"

//EventList is JSON struct of a list of events
type EventList struct {
	List []EventProperties `json:"events"`
}

//Event is JSOn struct of a single firewall's event
type Event struct {
	Properties EventProperties `json:"event"`
}

//EventProperties is JSON struct of an event properties
type EventProperties struct {
	ObjectType    string `json:"object_type"`
	RequestUUID   string `json:"request_uuid"`
	ObjectUUID    string `json:"object_uuid"`
	Activity      string `json:"activity"`
	RequestType   string `json:"request_type"`
	RequestStatus string `json:"request_status"`
	Change        string `json:"change"`
	Timestamp     string `json:"timestamp"`
	UserUUID      string `json:"user_uuid"`
}

//GetEventList gets a list of events
func (c *Client) GetEventList() ([]Event, error) {
	r := Request{
		uri:    apiEventBase,
		method: http.MethodGet,
	}
	var response EventList
	var events []Event
	err := r.execute(*c, &response)
	for _, properties := range response.List {
		events = append(events, Event{Properties: properties})
	}
	return events, err
}
