package service_query

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"time"
)

const timeoutCheckExistence = 1000 * time.Millisecond

//IsObjectExist checks if object exists
func IsObjectExist(client *gsclient.Client, service gsService, id string) (bool, error) {
	var isExist bool
	var err error
	switch service {
	case LoadbalancerService:
		_, err = client.GetLoadBalancer(id)
	case IPService:
		_, err = client.GetIP(id)
	case NetworkService:
		_, err = client.GetNetwork(id)
	case ServerService:
		_, err = client.GetServer(id)
	case SSHKeyService:
		_, err = client.GetSshkey(id)
	case StorageService:
		_, err = client.GetStorage(id)
	default:
		return isExist, fmt.Errorf("invalid service")
	}
	if err == nil {
		isExist = true
		return isExist, nil
	}
	//404 means this object does not exist
	//just return false and nil error
	if requestError, ok := err.(gsclient.RequestError); ok {
		if requestError.StatusCode == 404 {
			return isExist, nil
		}
	}
	return isExist, err
}
