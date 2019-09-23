package gsclient

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

//Config config for client
type Config struct {
	APIUrl     string
	UserUUID   string
	APIToken   string
	HTTPClient *http.Client
	logger     logrus.Logger
}

//NewConfiguration creates a new config
func NewConfiguration(apiURL string, uuid string, token string, debugMode bool) *Config {
	logLevel := logrus.InfoLevel
	if debugMode {
		logLevel = logrus.DebugLevel
	}

	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: logLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
		},
	}

	cfg := &Config{
		APIUrl:     apiURL,
		UserUUID:   uuid,
		APIToken:   token,
		HTTPClient: http.DefaultClient,
		logger:     logger,
	}
	return cfg
}
