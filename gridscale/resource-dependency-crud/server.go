package resource_dependency_crud

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/schema"
	service_query "github.com/terraform-providers/terraform-provider-gridscale/gridscale/service-query"
)

type serverRelation string

const (
	serverIsoimage = "server-isoimage"
	serverIP       = "server-ip"
	serverNetwork  = "server-network"
	serverStorage  = "server-storage"
)

//ServerDependencyClient is an wrapper of gsclient which is used for
//CRUD dependency of a server in gridscale's terraform provider
type ServerDependencyClient struct {
	gsc  *gsclient.Client
	Data *schema.ResourceData
}

//NewServerDepClient creates a new instance DependencyClient
func NewServerDepClient(gsc *gsclient.Client, d *schema.ResourceData) *ServerDependencyClient {
	return &ServerDependencyClient{gsc, d}
}

func (c *ServerDependencyClient) GetGSClient() *gsclient.Client {
	return c.gsc
}

//LinkStorages links a boot storage to a server
func (c *ServerDependencyClient) LinkStorages() error {
	d := c.Data
	client := c.GetGSClient()
	//Boot storage has to be attached first
	if attr, ok := d.GetOk("storage"); ok {
		for _, value := range attr.(*schema.Set).List() {
			storage := value.(map[string]interface{})
			isRelExist, _ := c.isServerRelationExists(serverStorage, storage["object_uuid"].(string))
			if storage["bootdevice"].(bool) && !isRelExist {
				err := client.LinkStorage(d.Id(), storage["object_uuid"].(string), storage["bootdevice"].(bool))
				if err != nil {
					return fmt.Errorf(
						"Error waiting for storage (%s) to be attached to server (%s): %s",
						storage["object_uuid"].(string),
						d.Id(),
						err,
					)
				}
			}
		}
	}
	//Attach additional storages
	if attr, ok := d.GetOk("storage"); ok {
		for _, value := range attr.(*schema.Set).List() {
			storage := value.(map[string]interface{})
			isRelExist, _ := c.isServerRelationExists(serverStorage, storage["object_uuid"].(string))
			if !storage["bootdevice"].(bool) && !isRelExist {
				err := client.LinkStorage(d.Id(), storage["object_uuid"].(string), storage["bootdevice"].(bool))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//LinkIPv4 links IPv4 address to a server
func (c *ServerDependencyClient) LinkIPv4() error {
	d := c.Data
	client := c.GetGSClient()
	if attr, ok := d.GetOk("ipv4"); ok {
		isRelExist, _ := c.isServerRelationExists(serverIP, attr.(string))
		if isRelExist {
			return nil
		}
		if client.GetIPVersion(attr.(string)) != 4 {
			return fmt.Errorf("The IP address with UUID %v is not version 4", attr.(string))
		}
		err := client.LinkIP(d.Id(), attr.(string))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for IP address (%s) to be attached to server (%s): %s",
				attr.(string),
				d.Id(),
				err,
			)
		}
	}
	return nil
}

//LinkIPv6 link an IPv6 address to a server
func (c *ServerDependencyClient) LinkIPv6() error {
	d := c.Data
	client := c.GetGSClient()
	if attr, ok := d.GetOk("ipv6"); ok {
		isRelExist, _ := c.isServerRelationExists(serverIP, attr.(string))
		if isRelExist {
			return nil
		}
		if client.GetIPVersion(attr.(string)) != 6 {
			return fmt.Errorf("The IP address with UUID %v is not version 6", attr.(string))
		}
		err := client.LinkIP(d.Id(), attr.(string))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for IP address (%s) to be attached to server (%s): %s",
				attr.(string),
				d.Id(),
				err,
			)
		}
	}
	return nil
}

//LinkISOImage links an ISO-image to a server
func (c *ServerDependencyClient) LinkISOImage() error {
	d := c.Data
	client := c.GetGSClient()
	if attr, ok := d.GetOk("isoimage"); ok {
		isRelExist, _ := c.isServerRelationExists(serverIsoimage, attr.(string))
		if isRelExist {
			return nil
		}
		err := client.LinkIsoImage(d.Id(), attr.(string))
		if err != nil {
			return fmt.Errorf(
				"Error waiting for iso-image (%s) to be attached to server (%s): %s",
				attr.(string),
				d.Id(),
				err,
			)
		}
	}
	return nil
}

//LinkNetworks links networks to server
func (c *ServerDependencyClient) LinkNetworks(isPublic bool) error {
	d := c.Data
	client := c.GetGSClient()
	if isPublic {
		publicNetwork, err := client.GetNetworkPublic()
		if err != nil {
			return err
		}
		isRelExist, _ := c.isServerRelationExists(serverNetwork, publicNetwork.Properties.ObjectUUID)
		if isRelExist {
			return nil
		}
		err = client.LinkNetwork(
			d.Id(),
			publicNetwork.Properties.ObjectUUID,
			"",
			false,
			0,
			nil,
			gsclient.FirewallRules{},
		)
		if err != nil {
			return fmt.Errorf(
				"Error waiting for public network (%s) to be attached to server (%s): %s",
				publicNetwork.Properties.ObjectUUID,
				d.Id(),
				err,
			)
		}
		return nil
	}
	if attr, ok := d.GetOk("network"); ok {
		for _, value := range attr.(*schema.Set).List() {
			network := value.(map[string]interface{})
			isRelExist, _ := c.isServerRelationExists(serverNetwork, network["object_uuid"].(string))
			if isRelExist {
				return nil
			}
			err := client.LinkNetwork(
				d.Id(),
				network["object_uuid"].(string),
				"",
				network["bootdevice"].(bool),
				0,
				nil,
				gsclient.FirewallRules{},
			)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for network (%s) to be attached to server (%s): %s",
					network["object_uuid"],
					d.Id(),
					err,
				)
			}
		}
	}
	return nil
}

//IsShutdownRequired checks if server is needed to be shutdown when updating
func (c *ServerDependencyClient) IsShutdownRequired() bool {
	var shutdownRequired bool
	d := c.Data
	if d.HasChange("cores") {
		old, new := d.GetChange("cores")
		if new.(int) < old.(int) || d.Get("legacy").(bool) { //Legacy systems don't support updating the memory while running
			shutdownRequired = true
		}
	}
	if d.HasChange("memory") {
		old, new := d.GetChange("memory")
		if new.(int) < old.(int) || d.Get("legacy").(bool) { //Legacy systems don't support updating the memory while running
			shutdownRequired = true
		}
	}
	if d.HasChange("ipv4") || d.HasChange("ipv6") || d.HasChange("storage") || d.HasChange("network") {
		shutdownRequired = true
	}
	return shutdownRequired
}

//UpdateISOImageRel updates relationship between a server and an ISO-image
func (c *ServerDependencyClient) UpdateISOImageRel() error {
	d := c.Data
	client := c.GetGSClient()
	var err error
	if d.HasChange("isoimage") {
		oldIso, newIso := d.GetChange("isoimage")
		if newIso == "" {
			isRelExist, _ := c.isServerRelationExists(serverIsoimage, oldIso.(string))
			if !isRelExist {
				return nil
			}
			err = client.UnlinkIsoImage(d.Id(), oldIso.(string))

		} else {
			isRelExist, _ := c.isServerRelationExists(serverIsoimage, newIso.(string))
			if isRelExist {
				return nil
			}
			err = client.LinkIsoImage(d.Id(), newIso.(string))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

//UpdateIPv4Rel updates relationship between a server and an IPv4 address
func (c *ServerDependencyClient) UpdateIPv4Rel() (bool, error) {
	needsPublicNetwork := true
	d := c.Data
	client := c.GetGSClient()
	var err error
	if d.HasChange("ipv4") {
		oldIp, newIp := d.GetChange("ipv4")
		if newIp == "" {
			isRelExist, _ := c.isServerRelationExists(serverIP, oldIp.(string))
			if isRelExist {
				err = client.UnlinkIP(d.Id(), oldIp.(string))
			}
		} else {
			isRelExist, _ := c.isServerRelationExists(serverIP, newIp.(string))
			if !isRelExist {
				err = client.LinkIP(d.Id(), newIp.(string))
			}
		}
		if err != nil {
			return needsPublicNetwork, err
		}
		if oldIp != "" {
			needsPublicNetwork = false
		}
	}
	return needsPublicNetwork, err
}

//UpdateIPv6Rel updates realtionship between a server and an IPv6 address
func (c *ServerDependencyClient) UpdateIPv6Rel() (bool, error) {
	needsPublicNetwork := true
	d := c.Data
	client := c.GetGSClient()
	var err error
	if d.HasChange("ipv6") {
		oldIp, newIp := d.GetChange("ipv6")
		if newIp == "" {
			isRelExist, _ := c.isServerRelationExists(serverIP, oldIp.(string))
			if isRelExist {
				err = client.UnlinkIP(d.Id(), oldIp.(string))
			}
		} else {
			isRelExist, _ := c.isServerRelationExists(serverIP, newIp.(string))
			if !isRelExist {
				err = client.LinkIP(d.Id(), newIp.(string))
			}
		}
		if err != nil {
			return needsPublicNetwork, err
		}
		if oldIp != "" {
			needsPublicNetwork = false
		}
	}
	return needsPublicNetwork, err
}

//UpdatePublicNetworkRel updates relationship between a server and a oublic network
func (c *ServerDependencyClient) UpdatePublicNetworkRel(isToLink bool) error {
	d := c.Data
	client := c.GetGSClient()
	publicNetwork, err := client.GetNetworkPublic()
	if err != nil {
		return err
	}
	isRelExist, _ := c.isServerRelationExists(serverNetwork, publicNetwork.Properties.ObjectUUID)
	if isToLink {
		if !isRelExist {
			err = client.LinkNetwork(d.Id(), publicNetwork.Properties.ObjectUUID, "", false, 0, []string{}, gsclient.FirewallRules{})
			if err != nil {
				return err
			}
		}
	} else {
		if isRelExist {
			err = client.UnlinkNetwork(d.Id(), publicNetwork.Properties.ObjectUUID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//UpdateOtherNetworkRel updates relationship between a server and networks
func (c *ServerDependencyClient) UpdateOtherNetworkRel() error {
	d := c.Data
	client := c.GetGSClient()
	var err error
	//It currently unlinks and relinks all networks if any network has changed. This could probably be done better, but this way is easy and works well
	if d.HasChange("network") {
		oldNetworks, newNetworks := d.GetChange("network")
		for _, value := range oldNetworks.(*schema.Set).List() {
			network := value.(map[string]interface{})
			isRelExist, _ := c.isServerRelationExists(serverNetwork, network["object_uuid"].(string))
			if network["object_uuid"].(string) != "" && isRelExist {
				err = client.UnlinkNetwork(d.Id(), network["object_uuid"].(string))
				if err != nil {
					return err
				}
			}
		}
		for _, value := range newNetworks.(*schema.Set).List() {
			network := value.(map[string]interface{})
			isRelExist, _ := c.isServerRelationExists(serverNetwork, network["object_uuid"].(string))
			if network["object_uuid"].(string) != "" && !isRelExist {
				err = client.LinkNetwork(d.Id(), network["object_uuid"].(string), "", network["bootdevice"].(bool), 0, []string{}, gsclient.FirewallRules{})
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//UpdateStorageRel updates relationship between a server and storages
func (c *ServerDependencyClient) UpdateStorageRel() error {
	d := c.Data
	client := c.GetGSClient()
	var err error
	if d.HasChange("storage") {
		oldStorages, newStorages := d.GetChange("storage")
		//unlink old storages if needed
		for _, value := range oldStorages.(*schema.Set).List() {
			oldStorage := value.(map[string]interface{})
			unlink := true
			for _, value := range newStorages.(*schema.Set).List() {
				newStorage := value.(map[string]interface{})
				if oldStorage["object_uuid"].(string) == newStorage["object_uuid"].(string) {
					unlink = false
					break
				}
			}
			isRelExist, _ := c.isServerRelationExists(serverStorage, oldStorage["object_uuid"].(string))
			if unlink && isRelExist {
				err = client.UnlinkStorage(d.Id(), oldStorage["object_uuid"].(string))
				if err != nil {
					return err
				}
			}
		}

		//link new storages if needed
		for _, value := range newStorages.(*schema.Set).List() {
			newStorage := value.(map[string]interface{})
			link := true
			for _, value := range oldStorages.(*schema.Set).List() {
				oldStorage := value.(map[string]interface{})
				if oldStorage["object_uuid"].(string) == newStorage["object_uuid"].(string) {
					link = false
					break
				}
			}
			isRelExist, _ := c.isServerRelationExists(serverStorage, newStorage["object_uuid"].(string))
			if link && !isRelExist {
				err = client.LinkStorage(d.Id(), newStorage["object_uuid"].(string), newStorage["bootdevice"].(bool))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

//isServerRelationExists check if relationship between a server and another service exists
func (c *ServerDependencyClient) isServerRelationExists(rel serverRelation, objID string) (bool, error) {
	client := c.GetGSClient()
	serverID := c.Data.Id()
	var isExist bool
	var err error
	switch rel {
	case serverIsoimage:
		isObjExist, _ := service_query.IsObjectExist(client, service_query.ISOImageService, objID)
		if !isObjExist {
			return false, nil
		}
		_, err = client.GetServerIsoImage(serverID, objID)
	case serverIP:
		isObjExist, _ := service_query.IsObjectExist(client, service_query.IPService, objID)
		if !isObjExist {
			return false, nil
		}
		_, err = client.GetServerIP(serverID, objID)
	case serverNetwork:
		isObjExist, _ := service_query.IsObjectExist(client, service_query.NetworkService, objID)
		if !isObjExist {
			return false, nil
		}
		_, err = client.GetServerNetwork(serverID, objID)
	case serverStorage:
		isObjExist, _ := service_query.IsObjectExist(client, service_query.StorageService, objID)
		if !isObjExist {
			return false, nil
		}
		_, err = client.GetServerStorage(serverID, objID)
	default:
		return isExist, fmt.Errorf("Invalid relationship")
	}
	if err == nil {
		isExist = true
		return isExist, err
	}
	//404 means this relationship does not exist
	//just return false and nil error
	if requestError, ok := err.(gsclient.RequestError); ok {
		if requestError.StatusCode == 404 {
			return isExist, nil
		}
	}
	return isExist, err
}
