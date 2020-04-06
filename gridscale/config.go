package gridscale

import (
	"context"
	"log"
	"os"

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
const (
	defaultAPIURL                    = "https://api.gridscale.io"
	defaultGSCTimeoutSecs            = 120
	defaultGSCDelayIntervalMilliSecs = 1000
	defaultGSCMaxNumberOfRetries     = 5
)

type Config struct {
	UserUUID    string
	APIToken    string
	APIUrl      string
	TimeoutSecs int
}

func (c *Config) Client() (*gsclient.Client, error) {
	var config *gsclient.Config
	// if api URL is configured, set the url in gsc
	apiURL := defaultAPIURL
	if c.APIUrl != "" {
		apiURL = c.APIUrl
	}

	//if timeout is configured, set the timeout in gsc
	timeoutSecs := defaultGSCTimeoutSecs
	if c.TimeoutSecs != 0 {
		timeoutSecs = c.TimeoutSecs
	}
	config = gsclient.NewConfiguration(
		apiURL,
		c.UserUUID,
		c.APIToken,
		os.Getenv("TF_LOG") != "",
		true,
		timeoutSecs,
		defaultGSCDelayIntervalMilliSecs,
		defaultGSCMaxNumberOfRetries,
	)

	client := gsclient.NewClient(config)

	log.Print("[INFO] gridscale client configured")

	//Make sure the credentials are correct by getting the server list
	//and init `globalServerStatusList` from fetched server list
	err := initGlobalServerStatusList(emptyCtx, client)

	return client, err
}
