package gridscale

import (
	"context"
	"log"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGridscaleIpv6() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleIpv6Create,
		Read:   resourceGridscaleIpRead,
		Delete: resourceGridscaleIpDelete,
		Update: resourceGridscaleIpUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Description: "Defines the IP address.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Optional:    true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "The location this object is placed.",
				Computed:    true,
			},
			"failover": {
				Type:        schema.TypeBool,
				Description: "Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.",
				Optional:    true,
				Default:     false,
			},
			"reverse_dns": {
				Type:        schema.TypeString,
				Description: "Defines the reverse DNS entry for the IP address (PTR Resource Record).",
				Optional:    true,
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
			"delete_block": {
				Type:        schema.TypeBool,
				Description: "Defines if the object is administratively blocked. If true, it can not be deleted by the user.",
				Computed:    true,
			},
			"usage_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "Defines the price for the current period since the last bill.",
				Computed:    true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleIpv6Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.IPCreateRequest{
		Family:     gsclient.IPv6Type,
		Name:       d.Get("name").(string),
		Failover:   d.Get("failover").(bool),
		ReverseDNS: d.Get("reverse_dns").(string),
		Labels:     convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateIP(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new Ipv%v has been set to %v", requestBody.Family, response.ObjectUUID)

	return resourceGridscaleIpRead(d, meta)
}
