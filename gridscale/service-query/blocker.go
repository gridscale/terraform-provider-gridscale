package service_query

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"time"
)

const (
	delayFetchingStatus  = 500 * time.Millisecond
	timeoutCheckDeletion = 1000 * time.Millisecond
)

const (
	provivisoningStatus = "in-provisioning"
	activeStatus        = "active"
)

type gsService string

const (
	LoadbalancerService gsService = "loadbalancer"
	IPService           gsService = "IP"
	NetworkService      gsService = "network"
	ServerService       gsService = "server"
	SSHKeyService       gsService = "sshkey"
	StorageService      gsService = "storage"
)

//BlockProvisoning blocks until the object's state is not in provisioning anymore
func BlockProvisoning(client *gsclient.Client, service gsService, id string) error {
	switch service {
	case LoadbalancerService:
		lb, err := client.GetLoadBalancer(id)
		for lb.Properties.Status == provivisoningStatus {
			lb, err = client.GetLoadBalancer(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for loadbalancer (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus)
		}
	case IPService:
		ip, err := client.GetIP(id)
		for ip.Properties.Status == provivisoningStatus {
			ip, err = client.GetIP(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for IP (%s) to be fetched: %s", id, err)
			}
			time.Sleep(500 * time.Millisecond)
		}
	case NetworkService:
		net, err := client.GetNetwork(id)
		for net.Properties.Status == provivisoningStatus {
			net, err = client.GetNetwork(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for network (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus)
		}
	case ServerService:
		server, err := client.GetServer(id)
		for server.Properties.Status == provivisoningStatus {
			server, err = client.GetServer(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for server (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus)
		}
	case SSHKeyService:
		sshKey, err := client.GetSshkey(id)
		for sshKey.Properties.Status == provivisoningStatus {
			sshKey, err = client.GetSshkey(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for ssh-key (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus)
		}
	case StorageService:
		storage, err := client.GetStorage(id)
		for storage.Properties.Status == provivisoningStatus {
			storage, err = client.GetStorage(id)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for storage (%s) to be fetched: %s", id, err)
			}
			time.Sleep(delayFetchingStatus)
		}
	default:
		return fmt.Errorf("invalid service")
	}
	return nil
}

//BlockDeletion blocks until an object is deleted successfully
func BlockDeletion(client *gsclient.Client, service gsService, id string) error {
	timer := time.After(timeoutCheckDeletion)
	var err error
	for {
		select {
		case <-timer:
			if err != nil {
				return fmt.Errorf("Timeout reached when waiting for %v (%v) to be deleted. Latest error: %v", service, id, err)
			}
			return fmt.Errorf("Timeout reached when waiting for %v (%v) to be deleted. Object still exists!", service, id)
		default:
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
				return fmt.Errorf("invalid service")
			}
			if err != nil {
				if requestError, ok := err.(gsclient.RequestError); ok {
					if requestError.StatusCode == 404 {
						return nil
					}
				}
			}
		}
	}
}
