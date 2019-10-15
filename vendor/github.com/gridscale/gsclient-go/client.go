package gsclient

import (
	"context"
	"errors"
	"path"
	"strings"
)

const (
	requestBase          = "/requests/"
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

//waitForRequestCompleted allows to wait for a request to complete
func (c *Client) waitForRequestCompleted(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("'id' is required")
	}
	return retryWithTimeout(func() (bool, error) {
		r := Request{
			uri:    path.Join(requestBase, id),
			method: "GET",
		}
		var response RequestStatus
		err := r.execute(ctx, *c, &response)
		if err != nil {
			return false, err
		}
		if response[id].Status == requestDoneStatus {
			c.cfg.logger.Info("Done with creating")
			return false, nil
		}
		return true, nil
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitFor404Status waits until server returns 404 status code
func (c *Client) waitFor404Status(ctx context.Context, uri, method string) error {
	return retryWithTimeout(func() (bool, error) {
		r := Request{
			uri:          uri,
			method:       method,
			skipPrint404: true,
		}
		err := r.execute(ctx, *c, nil)
		if err != nil {
			if requestError, ok := err.(RequestError); ok {
				if requestError.StatusCode == 404 {
					return false, nil
				}
			}
			return false, err
		}
		return true, nil
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}

//waitFor200Status waits until server returns 200 (OK) status code
func (c *Client) waitFor200Status(ctx context.Context, uri, method string) error {
	return retryWithTimeout(func() (bool, error) {
		r := Request{
			uri:          uri,
			method:       method,
			skipPrint404: true,
		}
		err := r.execute(ctx, *c, nil)
		if err != nil {
			if requestError, ok := err.(RequestError); ok {
				if requestError.StatusCode == 404 {
					return true, nil
				}
			}
			return false, err
		}
		return false, nil
	}, c.cfg.requestCheckTimeoutSecs, c.cfg.delayInterval)
}
