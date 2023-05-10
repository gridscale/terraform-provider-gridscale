package relationmanager

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	fwu "github.com/terraform-providers/terraform-provider-gridscale/gridscale/firewall-utils"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"
)

// firewallRuleTypes defines all types of firewall rules
var firewallRuleTypes = []string{"rules_v4_in", "rules_v4_out", "rules_v6_in", "rules_v6_out"}

// ServerRelationManger is an wrapper of gsclient which is used for
// managing relations of a server in gridscale terraform provider
type ServerRelationManger struct {
	gsc  *gsclient.Client
	data *schema.ResourceData
}

// NewServerRelationManger creates a new instance ServerRelationManger
func NewServerRelationManger(gsc *gsclient.Client, d *schema.ResourceData) *ServerRelationManger {
	return &ServerRelationManger{gsc, d}
}

// getGSClient returns gsclient from server relation manager
func (c ServerRelationManger) getGSClient() *gsclient.Client {
	return c.gsc
}

// getData returns resource data from server relation manager
func (c ServerRelationManger) getData() *schema.ResourceData {
	return c.data
}

// LinkStorages links storages to a server
// **Note: The first storage in the list will be automatically set as the boot device
func (c *ServerRelationManger) LinkStorages(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	if attr, ok := d.GetOk("storage"); ok {
		for _, value := range attr.([]interface{}) {
			storage := value.(map[string]interface{})
			err := client.LinkStorage(ctx, d.Id(), storage["object_uuid"].(string), false)
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
	return nil
}

// LinkIPv4 links IPv4 address to a server
func (c *ServerRelationManger) LinkIPv4(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	if attr, ok := d.GetOk("ipv4"); ok {
		//Check IP version
		if client.GetIPVersion(ctx, attr.(string)) != 4 {
			return fmt.Errorf("The IP address with UUID %v is not version 4", attr.(string))
		}
		err := client.LinkIP(ctx, d.Id(), attr.(string))
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

// LinkIPv6 link an IPv6 address to a server
func (c *ServerRelationManger) LinkIPv6(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	if attr, ok := d.GetOk("ipv6"); ok {
		//Check IP version
		if client.GetIPVersion(ctx, attr.(string)) != 6 {
			return fmt.Errorf("The IP address with UUID %v is not version 6", attr.(string))
		}
		err := client.LinkIP(ctx, d.Id(), attr.(string))
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

// LinkISOImage links an ISO image to a server
func (c *ServerRelationManger) LinkISOImage(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	if attr, ok := d.GetOk("isoimage"); ok {
		err := client.LinkIsoImage(ctx, d.Id(), attr.(string))
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

// LinkNetworks links networks to server
func (c *ServerRelationManger) LinkNetworks(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	if attrNetRel, ok := d.GetOk("network"); ok {
		for idx, value := range attrNetRel.([]interface{}) {
			// customFwRulesPtr is nil initially, that mean the fw is inactive
			var customFwRulesPtr *gsclient.FirewallRules
			network := value.(map[string]interface{})
			//Read custom firewall rules from `network` property (field)
			customFwRules := readCustomFirewallRules(network)
			// if customFwRules is not empty, customFwRulesPtr is not nil (fw is active)
			if !reflect.DeepEqual(customFwRules, gsclient.FirewallRules{}) {
				customFwRulesPtr = &customFwRules
			}
			err := client.LinkNetwork(
				ctx,
				d.Id(),
				network["object_uuid"].(string),
				network["firewall_template_uuid"].(string),
				network["bootdevice"].(bool),
				idx,
				nil,
				customFwRulesPtr,
			)
			if err != nil {
				return fmt.Errorf(
					"Error waiting for network (%s) to be attached to server (%s): %s",
					network["object_uuid"],
					d.Id(),
					err,
				)
			}

			if network["ip"].(string) != "" {
				// Assign DHCP IP to the server (if applicable).
				if err := client.UpdateNetworkPinnedServer(
					ctx,
					network["object_uuid"].(string),
					d.Id(),
					gsclient.PinServerRequest{
						IP: network["ip"].(string),
					},
				); err != nil {
					return fmt.Errorf(
						"Error waiting for assigning DHCP IP (%s) to server (%s) in network (%s): %s",
						network["ip"].(string),
						d.Id(),
						network["object_uuid"],
						err,
					)
				}
			}

		}
	}
	return nil
}

// readCustomFirewallRules reads custom firewall rules from a specific network
// returns `gsclient.FirewallRules` type variable
func readCustomFirewallRules(netData map[string]interface{}) gsclient.FirewallRules {
	//Init firewall rule variable
	var fwRules gsclient.FirewallRules

	//Loop through all firewall rule types
	//there are 4 types: "rules_v4_in", "rules_v4_out", "rules_v6_in", "rules_v6_out".
	for _, ruleType := range firewallRuleTypes {
		//Init array of firewall rules
		var rules []gsclient.FirewallRuleProperties
		//Check if the firewall rule type is declared in the current network
		if rulesInTypeAttr, ok := netData[ruleType]; ok {
			//Loop through all rules in the current firewall type
			for _, rulesInType := range rulesInTypeAttr.([]interface{}) {
				ruleProps := rulesInType.(map[string]interface{})
				ruleProperties := gsclient.FirewallRuleProperties{
					DstPort: ruleProps["dst_port"].(string),
					SrcPort: ruleProps["src_port"].(string),
					SrcCidr: ruleProps["src_cidr"].(string),
					Action:  ruleProps["action"].(string),
					Comment: ruleProps["comment"].(string),
					DstCidr: ruleProps["dst_cidr"].(string),
					Order:   ruleProps["order"].(int),
				}
				if ruleProps["protocol"].(string) == "tcp" {
					ruleProperties.Protocol = gsclient.TCPTransport
				} else if ruleProps["protocol"].(string) == "udp" {
					ruleProperties.Protocol = gsclient.UDPTransport
				}
				//Add rule to the array of rules
				rules = append(rules, ruleProperties)
			}
		}

		//Based on rule type to place the rules in the right property of fwRules variable
		if ruleType == "rules_v4_in" {
			fwRules.RulesV4In = fwu.AddDefaultFirewallInboundRules(rules, false) // add default rules
		} else if ruleType == "rules_v4_out" {
			fwRules.RulesV4Out = rules
		} else if ruleType == "rules_v6_in" {
			fwRules.RulesV6In = fwu.AddDefaultFirewallInboundRules(rules, true) // add default rules
		} else if ruleType == "rules_v6_out" {
			fwRules.RulesV6Out = rules
		}
	}
	return fwRules
}

// IsShutdownRequired checks if server is needed to be shutdown when updating
func (c *ServerRelationManger) IsShutdownRequired(ctx context.Context) bool {
	var shutdownRequired bool
	d := c.getData()
	hasServerNetListChange := c.hasServerNetworkListChanged(ctx)
	if d.HasChanges("cores", "memory", "ipv4", "ipv6", "hardware_profile", "hardware_profile_config", "auto_recovery", "user_data_base64") || hasServerNetListChange {
		shutdownRequired = true
	}
	return shutdownRequired
}

// UpdateISOImageRel updates relationship between a server and an ISO image
func (c *ServerRelationManger) UpdateISOImageRel(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	var err error
	//Check if ISO image field is changed
	if d.HasChange("isoimage") {
		oldIso, _ := d.GetChange("isoimage")
		//If there is an ISO image already linked to the server
		//Unlink it
		if oldIso != "" {
			//If 404 or 409, that means ISO image is already deleted => the relation between ISO image and server is deleted automatically
			err = errHandler.SuppressHTTPErrorCodes(
				client.UnlinkIsoImage(ctx, d.Id(), oldIso.(string)),
				http.StatusConflict,
				http.StatusNotFound,
			)
			if err != nil {
				return err
			}
		}
		//Link new ISO image (if there is one)
		err = c.LinkISOImage(ctx)
	}
	return err
}

// UpdateIPv4Rel updates relationship between a server and an IPv4 address
func (c *ServerRelationManger) UpdateIPv4Rel(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	var err error
	//If IPv4 field is changed
	if d.HasChange("ipv4") {
		oldIp, _ := d.GetChange("ipv4")
		//If there is an IPv4 address already linked to the server
		//Unlink it
		if oldIp != "" {
			//If 404 or 409, that means IP is already deleted => the relation between IP and server is deleted automatically
			err = errHandler.SuppressHTTPErrorCodes(
				client.UnlinkIP(ctx, d.Id(), oldIp.(string)),
				http.StatusConflict,
				http.StatusNotFound,
			)
			if err != nil {
				return err
			}
		}
		//Link new IPv4 (if there is one)
		err = c.LinkIPv4(ctx)
	}
	return err
}

// UpdateIPv6Rel updates relationship between a server and an IPv6 address
func (c *ServerRelationManger) UpdateIPv6Rel(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	var err error
	if d.HasChange("ipv6") {
		oldIp, _ := d.GetChange("ipv6")
		//If there is an IPv6 address already linked to the server
		//Unlink it
		if oldIp != "" {
			//If 404 or 409, that means IP is already deleted => the relation between IP and server is deleted automatically
			err = errHandler.SuppressHTTPErrorCodes(
				client.UnlinkIP(ctx, d.Id(), oldIp.(string)),
				http.StatusConflict,
				http.StatusNotFound,
			)
			if err != nil {
				return err
			}
		}
		//Link new IPv6 (if there is one)
		err = c.LinkIPv6(ctx)
	}
	return err
}

// RelinkAllNetworks relinks networks to the server. If there are no changes in
// server-network links, return without doing anything.
// Note: This action requires the server to be off.
func (c *ServerRelationManger) RelinkAllNetworks(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	// If there is no changes in server-net relations, return.
	if !d.HasChange("network") {
		return nil
	}
	//Unlink all old networks if there are any networks linked to the server
	oldNetworks, _ := d.GetChange("network")
	for _, value := range oldNetworks.([]interface{}) {
		network := value.(map[string]interface{})
		if network["object_uuid"].(string) != "" {
			//If 404 or 409, that means network is already deleted => the relation between network and server is deleted automatically
			err := errHandler.SuppressHTTPErrorCodes(
				client.UnlinkNetwork(ctx, d.Id(), network["object_uuid"].(string)),
				http.StatusConflict,
				http.StatusNotFound,
			)
			if err != nil {
				return err
			}
		}
	}
	//Links all new networks (if there are some)
	return c.LinkNetworks(ctx)
}

// UpdateNetRelsProperties update the properties of the server-network relations.
// If no changes, return without doing anything.
func (c *ServerRelationManger) UpdateNetRelsProperties(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	// If there is no changes in server-net relations, return.
	if !d.HasChange("network") {
		return nil
	}
	networkListIntf := d.Get("network").([]interface{})
	for idx, networkIntf := range networkListIntf {
		network := networkIntf.(map[string]interface{})
		//Read custom firewall rules from `network` property (field)
		customFwRules := readCustomFirewallRules(network)
		err := client.UpdateServerNetwork(
			ctx,
			d.Id(),
			network["object_uuid"].(string),
			gsclient.ServerNetworkRelationUpdateRequest{
				Ordering:             idx,
				BootDevice:           network["bootdevice"].(bool),
				Firewall:             &customFwRules,
				FirewallTemplateUUID: network["firewall_template_uuid"].(string),
			},
		)
		if err != nil {
			return fmt.Errorf(
				"error waiting for server(%s)-network(%s) relation to be updated: %s",
				d.Id(),
				network["object_uuid"],
				err,
			)
		}

		// Update DHCP IP assignment.
		if d.HasChange(fmt.Sprintf("network.%d.ip", idx)) {
			if network["ip"].(string) != "" {
				// Assign DHCP IP to the server (if applicable).
				if err := client.UpdateNetworkPinnedServer(
					ctx,
					network["object_uuid"].(string),
					d.Id(),
					gsclient.PinServerRequest{
						IP: network["ip"].(string),
					},
				); err != nil {
					return fmt.Errorf(
						"Error waiting for assigning DHCP IP (%s) to server (%s) in network (%s): %s",
						network["ip"].(string),
						d.Id(),
						network["object_uuid"],
						err,
					)
				}
			} else {
				if err := errHandler.SuppressHTTPErrorCodes(
					client.DeleteNetworkPinnedServer(
						ctx,
						network["object_uuid"].(string),
						d.Id(),
					),
					http.StatusNotFound,
				); err != nil {
					return fmt.Errorf(
						"Error waiting for removing DHCP IP (%s) from server (%s) in network (%s): %s",
						network["ip"].(string),
						d.Id(),
						network["object_uuid"],
						err,
					)
				}
			}
		}
	}
	return nil
}

// UpdateStoragesRel updates relationship between a server and storages
func (c *ServerRelationManger) UpdateStoragesRel(ctx context.Context) error {
	d := c.getData()
	client := c.getGSClient()
	var err error
	if d.HasChange("storage") {
		oldStorages, _ := d.GetChange("storage")
		for _, value := range oldStorages.([]interface{}) {
			storage := value.(map[string]interface{})
			if storage["object_uuid"].(string) != "" {
				//If 404 or 409, that means storage is already deleted => the relation between storage and server is deleted automatically
				err = errHandler.SuppressHTTPErrorCodes(
					client.UnlinkStorage(ctx, d.Id(), storage["object_uuid"].(string)),
					http.StatusConflict,
					http.StatusNotFound,
				)
				if err != nil {
					return err
				}
			}
		}
		//Links all new storages (if there are some)
		err = c.LinkStorages(ctx)
	}
	return err
}

// hasServerNetworkListChanged checks if a new network is being attached/detached
// to/from the server, or network ordering is changed.
func (c *ServerRelationManger) hasServerNetworkListChanged(ctx context.Context) bool {
	d := c.getData()
	oldNetList, newNetList := d.GetChange("network")
	var oldNetUUIDList []string
	var newNetUUIDList []string
	for _, netIntf := range oldNetList.([]interface{}) {
		net := netIntf.(map[string]interface{})
		oldNetUUIDList = append(oldNetUUIDList, net["object_uuid"].(string))
	}
	for _, netIntf := range newNetList.([]interface{}) {
		net := netIntf.(map[string]interface{})
		newNetUUIDList = append(newNetUUIDList, net["object_uuid"].(string))
	}
	// check if length of network list has changed.
	if len(oldNetUUIDList) != len(newNetUUIDList) {
		return true
	}
	// check if the attached network list has changed.
	if !reflect.DeepEqual(oldNetUUIDList, newNetUUIDList) {
		return true
	}
	return false
}
