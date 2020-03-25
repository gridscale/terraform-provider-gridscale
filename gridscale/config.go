package gridscale

import (
	"context"
	"log"

	"github.com/gridscale/gsclient-go/v2"
)

//Arrays can't be constants in Go, but these will be used as constants
var hardwareProfiles = []string{"default", "legacy", "nested", "cisco_csr", "sophos_utm", "f5_bigip", "q35", "q35_nested"}
var storageTypes = []string{"storage", "storage_high", "storage_insane"}
var availabilityZones = []string{"a", "b", "c"}
var loadbalancerAlgs = []string{"roundrobin", "leastconn"}
var passwordTypes = []string{"plain", "crypt"}
var firewallActionTypes = []string{"accept", "drop"}
var firewallRuleProtocols = []string{"udp", "tcp"}
var emptyCtx = context.Background()

const timeLayout = "2006-01-02 15:04:05"

type Config struct {
	UserUUID        string
	ProjectAPIToken map[string]string
	APIUrl          string
}

func (c *Config) Clients() (map[string]*gsclient.Client, error) {
	projectClients := make(map[string]*gsclient.Client)
	for projectName, token := range c.ProjectAPIToken {
		config := gsclient.DefaultConfiguration(
			c.UserUUID,
			token,
		)
		log.Print("[INFO] gridscale client configured")
		projectClients[projectName] = gsclient.NewClient(config)
	}

	//Make sure the credentials are correct by getting the server list
	//and init `globalServerStatusList` from fetched server list
	err := initGlobalServerStatusList(emptyCtx, projectClients)
	return projectClients, err
}
