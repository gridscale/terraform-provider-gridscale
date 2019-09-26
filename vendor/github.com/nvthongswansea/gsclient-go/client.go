package gsclient

import (
	"context"
	"fmt"
	"path"
	"time"
)

const (
	apiServerBase        = "/objects/servers"
	apiStorageBase       = "/objects/storages"
	apiNetworkBase       = "/objects/networks"
	apiIPBase            = "/objects/ips"
	apiSshkeyBase        = "/objects/sshkeys"
	apiTemplateBase      = "/objects/templates"
	apiLoadBalancerBase  = "/objects/loadbalancers"
	apiPaaSBase          = "/objects/paas"
	apiISOBase           = "/objects/isoimages"
	apiObjectStorageBase = "/objects/objectstorages"
	apiFirewallBase      = "/objects/firewalls"
	apiLocationBase      = "/objects/locations"
	apiEventBase         = "/objects/events"
	apiLabelBase         = "/objects/labels"
	apiDeletedBase       = "/objects/deleted"
)

const (
	activeStatus      = "active"
	requestDoneStatus = "done"
)

//Client struct of a gridscale golang client
type Client struct {
	cfg *Config
}

//NewClient creates new gridscale golang client
func NewClient(c *Config) *Client {
	client := &Client{
		cfg: c,
	}
	return client
}

//waitForRequestCompleted allows to wait for a request to complete. Timeouts are currently hardcoded
func (c *Client) waitForRequestCompleted(ctx context.Context, id string) error {
	r := Request{
		uri:    path.Join("/requests/", id),
		method: "GET",
	}
	timer := time.After(c.cfg.requestCheckTimeoutSecs)
	delayInterval := c.cfg.delayInterval
	for {
		select {
		case <-timer:
			c.cfg.logger.Errorf("Timeout reached when waiting for request %v to complete", id)
			return fmt.Errorf("Timeout reached when waiting for request %v to complete", id)
		default:
			time.Sleep(delayInterval) //delay the request, so we don't do too many requests to the server
			var response RequestStatus
			r.execute(ctx, *c, &response)
			if response[id].Status == requestDoneStatus {
				c.cfg.logger.Info("Done with creating")
				return nil
			}
		}
	}
}
