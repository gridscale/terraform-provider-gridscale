package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
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
	errorPrefix := "read public network datasource -"
	network, err := client.GetNetworkPublic(context.Background())
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
