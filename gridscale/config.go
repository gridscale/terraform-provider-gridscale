package gridscale

import (
	"../gsclient"
	"log"
)

type Config struct {
	UserUUID string
	APIToken string
}

func (c *Config) Client() (*gsclient.Client, error) {
	config := gsclient.NewConfiguration(c.UserUUID, c.APIToken)
	client := gsclient.NewClient(config)

	log.Print("[INFO] gridscale client configured")

	//Make sure the credentials are correct by getting the server list
	_, err := client.GetServerList()

	return client, err
}
