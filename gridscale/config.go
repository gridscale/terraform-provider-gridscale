package gridscale

import (
	"github.com/gridscale/gsclient-go"
	"log"
	"net/http"
)

//Arrays can't be constants in Go, but these will be used as constants
var HardwareProfiles = []string{"default", "legacy", "nested", "cisco_csr", "sophos_utm", "f5_bigip", "q35"}
var StorageTypes = []string{"storage", "storage_high", "storage_insane"}
var AvailabilityZones = []string{"a", "b", "c"}

type Config struct {
	UserUUID string
	APIToken string
	APIUrl   string
}

func (c *Config) Client() (*gsclient.Client, error) {
	//config := gsclient.NewConfiguration(c.UserUUID, c.APIToken)
	config := gsclient.Config{
		APIUrl:     c.APIUrl,
		UserUUID:   c.UserUUID,
		APIToken:   c.APIToken,
		HTTPClient: http.DefaultClient,
	}
	client := gsclient.NewClient(&config)

	log.Print("[INFO] gridscale client configured")

	//Make sure the credentials are correct by getting the server list
	_, err := client.GetServerList()

	return client, err
}
