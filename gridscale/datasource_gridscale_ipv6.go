package gridscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/nvthongswansea/gsclient-go"
)

func dataSourceGridscaleIpv6() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleIpv6Read,

		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"ip": {
				Type:        schema.TypeString,
				Description: "Defines the IP Address.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to",
				Computed:    true,
			},
			"failover": {
				Type:        schema.TypeBool,
				Description: "Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.",
				Computed:    true,
			},
			"reverse_dns": {
				Type:        schema.TypeString,
				Description: "Defines the reverse DNS entry for the IP Address (PTR Resource Record).",
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change.",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"delete_block": {
				Type:        schema.TypeBool,
				Description: "Defines if the object is administratively blocked. If true, it can not be deleted by the user.",
				Computed:    true,
			},
			"usage_in_minutes": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "Defines the price for the current period since the last bill.",
				Computed:    true,
			},
		},
	}
}

func dataSourceGridscaleIpv6Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	ip, err := client.GetIP(emptyCtx, id)

	if err == nil {
		d.SetId(ip.Properties.ObjectUUID)
		d.Set("ip", ip.Properties.IP)
		d.Set("name", ip.Properties.Name)
		d.Set("prefix", ip.Properties.Prefix)
		d.Set("location_uuid", ip.Properties.LocationUUID)
		d.Set("failover", ip.Properties.Failover)
		d.Set("status", ip.Properties.Status)
		d.Set("reverse_dns", ip.Properties.ReverseDNS)
		d.Set("location_country", ip.Properties.LocationCountry)
		d.Set("location_iata", ip.Properties.LocationIata)
		d.Set("location_name", ip.Properties.LocationName)
		d.Set("create_time", ip.Properties.CreateTime)
		d.Set("change_time", ip.Properties.ChangeTime)
		d.Set("delete_block", ip.Properties.DeleteBlock)
		d.Set("usage_in_minutes", ip.Properties.UsagesInMinutes)
		d.Set("current_price", ip.Properties.CurrentPrice)

		if err = d.Set("labels", ip.Properties.Labels); err != nil {
			return fmt.Errorf("Error setting labels: %v", err)
		}
	}

	return err
}
