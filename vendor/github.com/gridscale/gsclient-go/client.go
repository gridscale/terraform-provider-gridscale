package gsclient

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
