package gridscale

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/gridscale/gsclient-go/v3"
	"gopkg.in/yaml.v2"
)

//Arrays can't be constants in Go, but these will be used as constants
var hardwareProfiles = []string{"default", "legacy", "nested", "cisco_csr", "sophos_utm", "f5_bigip", "q35", "q35_nested"}
var storageTypes = []string{"storage", "storage_high", "storage_insane"}
var storageVariants = []string{"distributed", "local"}
var availabilityZones = []string{"a", "b", "c"}
var loadbalancerAlgs = []string{"roundrobin", "leastconn"}
var passwordTypes = []string{"plain", "crypt"}
var firewallActionTypes = []string{"accept", "drop"}
var firewallRuleProtocols = []string{"udp", "tcp"}
var marketplaceAppCategories = []string{"CMS", "project management", "Adminpanel", "Collaboration", "Cloud Storage", "Archiving"}
var postgreSQLPerformanceClasses = []string{"standard", "high", "insane", "ultra"}
var msSQLServerPerformanceClasses = []string{"standard", "high", "insane", "ultra"}
var mariaDBPerformanceClasses = []string{"standard", "high", "insane", "ultra"}

const timeLayout = "2006-01-02 15:04:05"
const (
	defaultAPIURL                    = "https://api.gridscale.io"
	defaultGSCDelayIntervalMilliSecs = 1000
	defaultGSCMaxNumberOfRetries     = 1
)

const serverShutdownTimeoutSecs = 120

type Config struct {
	UserUUID    string
	APIToken    string
	APIUrl      string
	DelayIntMs  int
	MaxNRetries int
	HTTPHeaders map[string]string
}

func (c *Config) Client() (*gsclient.Client, error) {
	// if api URL is configured, set the url in gsc
	apiURL := defaultAPIURL
	if c.APIUrl != "" {
		apiURL = c.APIUrl
	}
	delayIntMs := defaultGSCDelayIntervalMilliSecs
	if c.DelayIntMs != 0 {
		delayIntMs = c.DelayIntMs
	}
	maxNRetries := defaultGSCMaxNumberOfRetries
	if c.MaxNRetries != 0 {
		maxNRetries = c.MaxNRetries
	}
	config := gsclient.NewConfiguration(
		apiURL,
		c.UserUUID,
		c.APIToken,
		os.Getenv("TF_LOG") != "",
		true,
		delayIntMs,
		maxNRetries,
	)

	client := gsclient.NewClient(config)
	//Add HTTP headers to gs client
	client.WithHTTPHeaders(c.HTTPHeaders)

	log.Print("[INFO] gridscale client configured")

	//Make sure the credentials are correct by getting the server list
	//and init `globalServerStatusList` from fetched server list
	err := initGlobalServerStatusList(context.Background(), client)

	return client, err
}

// GSCloudAccountEntry represents a single account in the config file of gscloud.
type GSCloudAccountEntry struct {
	Name   string `yaml:"name"`
	UserID string `yaml:"userId"`
	Token  string `yaml:"token"`
	URL    string `yaml:"url"`
}

// GSCloudConfig are all configuration settings parsed from a configuration file of gscloud.
type GSCloudConfig struct {
	Accounts []GSCloudAccountEntry `yaml:"accounts"`
}

func getGSCloudConfigFromPath(configFile string) (GSCloudConfig, error) {
	config := GSCloudConfig{}
	fileContent, err := os.ReadFile(configFile)
	if err != nil {
		return GSCloudConfig{}, err
	}
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		return GSCloudConfig{}, err
	}
	if len(config.Accounts) == 0 {
		return GSCloudConfig{}, errors.New("no configuration in the config file.")
	}
	return config, nil
}
