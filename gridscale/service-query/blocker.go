package service_query

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/resource"
	"time"
)

const (
	delayFetchingStatus  = 500 * time.Millisecond
	timeoutCheckDeletion = 1 * time.Minute
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
	ISOImageService     gsService = "isoimage"
)

//BlockProvisoning blocks until the object's state is not in provisioning anymore
func BlockProvisoning(client *gsclient.Client, service gsService, id string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		time.Sleep(delayFetchingStatus)
		switch service {
		case LoadbalancerService:
			lb, err := client.GetLoadBalancer(id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf(
					"Error waiting for loadbalancer (%s) to be fetched: %s", id, err))
			}
			if lb.Properties.Status != activeStatus {
				return resource.RetryableError(fmt.Errorf("Status of loadbalancer %s is not active", id))
			}
			return nil
		case IPService:
			ip, err := client.GetIP(id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf(
					"Error waiting for IP (%s) to be fetched: %s", id, err))
			}
			if ip.Properties.Status != activeStatus {
				return resource.RetryableError(fmt.Errorf("Status of IP %s is not active", id))
			}
			return nil
		case NetworkService:
			net, err := client.GetNetwork(id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf(
					"Error waiting for network (%s) to be fetched: %s", id, err))
			}
			if net.Properties.Status != activeStatus {
				return resource.RetryableError(fmt.Errorf("Status of network %s is not active", id))
			}
			return nil
		case ServerService:
			server, err := client.GetServer(id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf(
					"Error waiting for server (%s) to be fetched: %s", id, err))
			}
			if server.Properties.Status != activeStatus {
				return resource.RetryableError(fmt.Errorf("Status of server %s is not active", id))
			}
			return nil
		case SSHKeyService:
			sshKey, err := client.GetSshkey(id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf(
					"Error waiting for sshKey (%s) to be fetched: %s", id, err))
			}
			if sshKey.Properties.Status != activeStatus {
				return resource.RetryableError(fmt.Errorf("Status of sshKey %s is not active", id))
			}
			return nil
		case StorageService:
			storage, err := client.GetStorage(id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf(
					"Error waiting for storage (%s) to be fetched: %s", id, err))
			}
			if storage.Properties.Status != activeStatus {
				return resource.RetryableError(fmt.Errorf("Status of storage %s is not active", id))
			}
			return nil
		default:
			return resource.NonRetryableError(fmt.Errorf("invalid service"))
		}
	})
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
