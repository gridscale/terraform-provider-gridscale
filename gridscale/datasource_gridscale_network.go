package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGridscaleNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleNetworkRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Computed:    true,
			},
			"l2security": {
				Type:        schema.TypeBool,
				Description: "MAC spoofing protection - filters layer2 and ARP traffic based on source MAC",
				Computed:    true,
			},
			"dhcp_active": {
				Type:        schema.TypeBool,
				Description: "Enable DHCP.",
				Computed:    true,
			},
			"dhcp_range": {
				Type:        schema.TypeString,
				Description: "The general IP Range configured for this network (/24 for private networks). If it is not set, gridscale internal default range is used.",
				Computed:    true,
			},
			"dhcp_gateway": {
				Type:        schema.TypeString,
				Description: "The IP address reserved and communicated by the dhcp service to be the default gateway.",
				Computed:    true,
			},
			"dhcp_dns": {
				Type:        schema.TypeString,
				Description: "DHCP DNS.",
				Computed:    true,
			},
			"dhcp_reserved_subnet": {
				Type:        schema.TypeSet,
				Description: "Subrange within the IP range",
				Computed:    true,
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
				Description: "status indicates the status of the object",
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
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read network (%s) datasource-", id)

	network, err := client.GetNetwork(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(network.Properties.ObjectUUID)
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
