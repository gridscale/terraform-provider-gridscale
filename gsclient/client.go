package gsclient

import (
	"fmt"
	"time"
)

type Client struct {
	cfg			*Config
}

type isReady func(string, chan bool)

func NewClient(c *Config) *Client {
	client := &Client{
		cfg:	c,
	}

	return client
}

func (c* Client) WaitForState(id string, isready isReady) error	{
	var inprogress chan bool
	select {
		case <- inprogress:
			 isready(id, inprogress)
		case <-time.After(5 * time.Second):
			return fmt.Errorf("Timeout reached when waiting for object %v to become active", id)
	}
	return nil
}