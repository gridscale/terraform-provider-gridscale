package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
)

func resourceGridscaleNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleNetworkCreate,
		Read:   resourceGridscaleNetworkRead,
		Delete: resourceGridscaleNetworkDelete,
		Update: resourceGridscaleNetworkUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
			},
			"l2security": {
				Type:        schema.TypeBool,
				Description: "MAC spoofing protection - filters layer2 and ARP traffic based on source MAC",
				Optional:    true,
				Default:     false,
			},
			"dhcp_active": {
				Type:        schema.TypeBool,
				Description: "Enable DHCP.",
				Optional:    true,
			},
			"dhcp_range": {
				Type:        schema.TypeString,
				Description: "The general IP Range configured for this network (/24 for private networks). If it is not set, gridscale internal default range is used.",
				Optional:    true,
			},
			"dhcp_gateway": {
				Type:        schema.TypeString,
				Description: "The IP address reserved and communicated by the dhcp service to be the default gateway.",
				Optional:    true,
			},
			"dhcp_dns": {
				Type:        schema.TypeString,
				Description: "DHCP DNS.",
				Optional:    true,
			},
			"dhcp_reserved_subnet": {
				Type:        schema.TypeSet,
				Description: "Subrange within the IP range",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"auto_assigned_servers": {
				Type:        schema.TypeSet,
				Description: "Contains IP addresses of all servers in the network which got a designated IP by the DHCP server.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"pinned_servers": {
				Type:        schema.TypeSet,
				Description: "Contains IP addresses of all servers in the network which got a designated IP by the user.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object.",
				Computed:    true,
			},
			"network_type": {
				Type:        schema.TypeString,
				Description: "The type of this network, can be mpls, breakout or network.",
				Computed:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "The location this object is placed.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Two digit country code (ISO 3166-2) of the location where this object is placed.",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters",
				Computed:    true,
			},
			"delete_block": {
				Type:        schema.TypeBool,
				Description: "If deleting this network is allowed",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read network (%s) resource -", d.Id())
	network, err := client.GetNetwork(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	if err = d.Set("name", network.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", network.Properties.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("l2security", network.Properties.L2Security); err != nil {
		return fmt.Errorf("%s error setting l2security: %v", errorPrefix, err)
	}
	if err = d.Set("dhcp_active", network.Properties.DHCPActive); err != nil {
		return fmt.Errorf("%s error setting dhcp_active: %v", errorPrefix, err)
	}
	if err = d.Set("dhcp_range", network.Properties.DHCPRange); err != nil {
		return fmt.Errorf("%s error setting dhcp_range: %v", errorPrefix, err)
	}
	if err = d.Set("dhcp_gateway", network.Properties.DHCPGateway); err != nil {
		return fmt.Errorf("%s error setting dhcp_gateway: %v", errorPrefix, err)
	}
	if err = d.Set("dhcp_dns", network.Properties.DHCPDNS); err != nil {
		return fmt.Errorf("%s error setting dhcp_dns: %v", errorPrefix, err)
	}
	if err = d.Set("dhcp_reserved_subnet", network.Properties.DHCPReservedSubnet); err != nil {
		return fmt.Errorf("%s error setting dhcp_reserved_subnet: %v", errorPrefix, err)
	}

	autoAssignedServers := make([]interface{}, 0)
	for _, value := range network.Properties.AutoAssignedServers {
		serverWIP := map[string]interface{}{
			"server_uuid": value.ServerUUID,
			"ip":          value.IP,
		}
		autoAssignedServers = append(autoAssignedServers, serverWIP)
	}
	if err = d.Set("auto_assigned_servers", autoAssignedServers); err != nil {
		return fmt.Errorf("%s error setting auto_assigned_servers: %v", errorPrefix, err)
	}

	pinnedServers := make([]interface{}, 0)
	for _, value := range network.Properties.PinnedServers {
		serverWIP := map[string]interface{}{
			"server_uuid": value.ServerUUID,
			"ip":          value.IP,
		}
		pinnedServers = append(pinnedServers, serverWIP)
	}
	if err = d.Set("pinned_servers", pinnedServers); err != nil {
		return fmt.Errorf("%s error setting pinned_servers: %v", errorPrefix, err)
	}

	if err = d.Set("status", network.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("network_type", network.Properties.NetworkType); err != nil {
		return fmt.Errorf("%s error setting network_type: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", network.Properties.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", network.Properties.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", network.Properties.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("delete_block", network.Properties.DeleteBlock); err != nil {
		return fmt.Errorf("%s error setting delete_block: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", network.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", network.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}

	if err = d.Set("labels", network.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	return nil
}

func resourceGridscaleNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update network (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.NetworkUpdateRequest{
		Name:       d.Get("name").(string),
		L2Security: d.Get("l2security").(bool),
		Labels:     &labels,
	}
	if d.HasChange("dhcp_active") {
		dhcpActive := d.Get("dhcp_active").(bool)
		requestBody.DHCPActive = &dhcpActive
	}
	if d.HasChange("dhcp_range") {
		dhcpRange := d.Get("dhcp_range").(string)
		requestBody.DHCPRange = &dhcpRange
	}
	if d.HasChange("dhcp_gateway") {
		dhcpGateway := d.Get("dhcp_gateway").(string)
		requestBody.DHCPGateway = &dhcpGateway
	}
	if d.HasChange("dhcp_dns") {
		dhcpDNS := d.Get("dhcp_dns").(string)
		requestBody.DHCPDNS = &dhcpDNS
	}
	if d.HasChange("dhcp_reserved_subnet") {
		dhcpReservedSubnet := convSOStrings(d.Get("dhcp_reserved_subnet").(*schema.Set).List())
		requestBody.DHCPReservedSubnet = &dhcpReservedSubnet
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateNetwork(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.NetworkCreateRequest{
		Name:       d.Get("name").(string),
		L2Security: d.Get("l2security").(bool),
		Labels:     convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	if dhcpActiveIntf, ok := d.GetOk("dhcp_active"); ok {
		requestBody.DHCPActive = dhcpActiveIntf.(bool)
	}
	if dhcpRangeIntf, ok := d.GetOk("dhcp_range"); ok {
		requestBody.DHCPRange = dhcpRangeIntf.(string)
	}
	if dhcpGatewayIntf, ok := d.GetOk("dhcp_gateway"); ok {
		requestBody.DHCPGateway = dhcpGatewayIntf.(string)
	}
	if dhcpDNSIntf, ok := d.GetOk("dhcp_dns"); ok {
		requestBody.DHCPDNS = dhcpDNSIntf.(string)
	}
	if dhcpReservedSubnetIntf, ok := d.GetOk("dhcp_reserved_subnet"); ok {
		requestBody.DHCPReservedSubnet = convSOStrings(dhcpReservedSubnetIntf.(*schema.Set).List())
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateNetwork(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for network %v has been set to %v", requestBody.Name, response.ObjectUUID)

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete network (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	net, err := client.GetNetwork(ctx, d.Id())
	//In case of 404, don't catch the error
	if errHandler.RemoveErrorContainsHTTPCodes(err, http.StatusNotFound) != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	//Stop all servers relating to this network address if there is one
	for _, server := range net.Properties.Relations.Servers {
		unlinkNetAction := func(ctx context.Context) error {
			//No need to unlink when server returns 409 or 404
			err := errHandler.RemoveErrorContainsHTTPCodes(
				client.UnlinkNetwork(ctx, server.ObjectUUID, d.Id()),
				http.StatusConflict,
				http.StatusNotFound,
			)
			return err
		}
		//UnlinkNetwork requires the server to be off
		err = globalServerStatusList.runActionRequireServerOff(ctx, client, server.ObjectUUID, false, unlinkNetAction)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}

	err = errHandler.RemoveErrorContainsHTTPCodes(
		client.DeleteNetwork(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
