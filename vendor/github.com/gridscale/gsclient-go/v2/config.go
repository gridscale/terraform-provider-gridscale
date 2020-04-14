package gsclient

import (
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultCheckRequestTimeoutSecs = 120
	defaultMaxNumberOfRetries      = 5
	defaultDelayIntervalMilliSecs  = 1000
	version                        = "2.2.2"
	defaultAPIURL                  = "https://api.gridscale.io"
	resourceActiveStatus           = "active"
	requestDoneStatus              = "done"
	requestFailStatus              = "failed"
	bodyType                       = "application/json"
)

//Config config for client
type Config struct {
	apiURL                  string
	userUUID                string
	apiToken                string
	userAgent               string
	sync                    bool
	httpClient              *http.Client
	requestCheckTimeoutSecs time.Duration
	delayInterval           time.Duration
	maxNumberOfRetries      int
	logger                  logrus.Logger
}

//NewConfiguration creates a new config
//
//- Parameters:
//		+ apiURL string: base URL of API.
//		+ uuid string: UUID of user.
//		+ token string: API token.
//		+ debugMode bool: true => run client in debug mode.
//		+ sync bool: true => client is in synchronous mode. The client will block until Create/Update/Delete processes
//		are completely finished. It is safer to set this parameter to `true`.
//		+ requestCheckTimeoutSecs int: Timeout (in second) for checking requests (for synchronous feature)
//		+ delayIntervalMilliSecs int: delay (in MilliSecond) between requests when checking request (or retry 5xx, 424 error code)
//		+ maxNumberOfRetries int: number of retries when server returns 5xx, 424 error code.
func NewConfiguration(apiURL string, uuid string, token string, debugMode, sync bool, requestCheckTimeoutSecs,
	delayIntervalMilliSecs, maxNumberOfRetries int) *Config {
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
		apiURL:                  apiURL,
		userUUID:                uuid,
		apiToken:                token,
		userAgent:               "gsclient-go/" + version + " (" + runtime.GOOS + ")",
		sync:                    sync,
		httpClient:              http.DefaultClient,
		logger:                  logger,
		requestCheckTimeoutSecs: time.Duration(requestCheckTimeoutSecs) * time.Second,
		delayInterval:           time.Duration(delayIntervalMilliSecs) * time.Millisecond,
		maxNumberOfRetries:      maxNumberOfRetries,
	}
	return cfg
}

//DefaultConfiguration creates a default configuration
func DefaultConfiguration(uuid string, token string) *Config {
	logger := logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: false,
		},
	}
	cfg := &Config{
		apiURL:                  defaultAPIURL,
		userUUID:                uuid,
		apiToken:                token,
		userAgent:               "gsclient-go/" + version + " (" + runtime.GOOS + ")",
		sync:                    true,
		httpClient:              http.DefaultClient,
		logger:                  logger,
		requestCheckTimeoutSecs: time.Duration(defaultCheckRequestTimeoutSecs) * time.Second,
		delayInterval:           time.Duration(defaultDelayIntervalMilliSecs) * time.Millisecond,
		maxNumberOfRetries:      defaultMaxNumberOfRetries,
	}
	return cfg
}
