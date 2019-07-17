package gsclient

import "net/http"

type Config struct {
	APIUrl   string
	UserUUID string
	APIToken string

	HTTPClient *http.Client
}

func NewConfiguration(uuid string, token string) *Config {
	cfg := &Config{
		APIUrl:     "https://api.gridscale.io",
		UserUUID:   uuid,
		APIToken:   token,
		HTTPClient: http.DefaultClient,
	}
	return cfg
}
