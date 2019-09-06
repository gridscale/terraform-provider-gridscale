package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"time"
)

type gsService string

const delayFetchingStatus = 500

const (
	provivisoningStatus = "in-provisioning"
	activeStatus        = "active"
)

const (
	loadbalancerService gsService = "loadbalancer"
	ipService           gsService = "IP"
	networkService      gsService = "network"
	serverService       gsService = "server"
	sshKeyService       gsService = "sshkey"
	storageService      gsService = "storage"
)

type objectGetter func(id string) error

//convSOStrings converts slice of interfaces to slice of strings
func convSOStrings(interfaceList []interface{}) []string {
	var labels []string
	for _, labelInterface := range interfaceList {
		labels = append(labels, labelInterface.(string))
	}
	return labels
}

func pauseWhenProvisoning(client *gsclient.Client, service gsService, id string) error {
	switch service {
	case loadbalancerService:
		lb, err := client.GetLoadBalancer(id)
		for lb.Properties.Status == provivisoningStatus {
			lb, err = client.GetLoadBalancer(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for loadbalancer (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus * time.Millisecond)
		}
	case ipService:
		ip, err := client.GetIP(id)
		for ip.Properties.Status == provivisoningStatus {
			ip, err = client.GetIP(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for IP (%s) to be fetched: %s", id, err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	case networkService:
		net, err := client.GetNetwork(id)
		for net.Properties.Status == provivisoningStatus {
			net, err = client.GetNetwork(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for network (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus * time.Millisecond)
		}
	case serverService:
		server, err := client.GetServer(id)
		for server.Properties.Status == provivisoningStatus {
			server, err = client.GetServer(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for server (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus * time.Millisecond)
		}
	case sshKeyService:
		sshKey, err := client.GetSshkey(id)
		for sshKey.Properties.Status == provivisoningStatus {
			sshKey, err = client.GetSshkey(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for ssh-key (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus * time.Millisecond)
		}
	case storageService:
		storage, err := client.GetStorage(id)
		for storage.Properties.Status == provivisoningStatus {
			storage, err = client.GetStorage(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for storage (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus * time.Millisecond)
		}
	default:
		return fmt.Errorf("invalid service")
	}
	return nil
}
