package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGridscalePublicNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscalePublicNetworkRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"l2security": {
				Type:        schema.TypeBool,
				Description: "MAC spoofing protection - filters layer2 and ARP traffic based on source MAC",
				Computed:    true,
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
				Description: "Helps to identify which datacenter an object belongs to",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Formatted by the 2 digit country code (ISO 3166-2) of the host country",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters",
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

func dataSourceGridscalePublicNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	network, err := client.GetNetworkPublic(emptyCtx)

	if err == nil {
		d.SetId(network.Properties.ObjectUUID)
		d.Set("name", network.Properties.Name)
		d.Set("location_uuid", network.Properties.LocationUUID)
		d.Set("l2security", network.Properties.L2Security)
		d.Set("status", network.Properties.Status)
		d.Set("network_type", network.Properties.NetworkType)
		d.Set("location_country", network.Properties.LocationCountry)
		d.Set("location_iata", network.Properties.LocationIata)
		d.Set("location_name", network.Properties.LocationName)
		d.Set("delete_block", network.Properties.DeleteBlock)
		d.Set("create_time", network.Properties.CreateTime.String())
		d.Set("change_time", network.Properties.ChangeTime.String())
		if err = d.Set("labels", network.Properties.Labels); err != nil {
			return fmt.Errorf("Error setting labels: %v", err)
		}
	}

	return err
}
