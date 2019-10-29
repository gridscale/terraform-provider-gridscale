package gridscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gridscale/gsclient-go"
)

func resourceGridscaleIpv4() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleIpv4Create,
		Read:   resourceGridscaleIpRead,
		Delete: resourceGridscaleIpDelete,
		Update: resourceGridscaleIpUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Description: "Defines the IP Address.",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Optional:    true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"failover": {
				Type:        schema.TypeBool,
				Description: "Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.",
				Optional:    true,
				Default:     false,
			},
			"reverse_dns": {
				Type:        schema.TypeString,
				Description: "Defines the reverse DNS entry for the IP Address (PTR Resource Record).",
				Optional:    true,
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
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(time.Minute * 3),
		},
	}
}

func resourceGridscaleIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	ip, err := client.GetIP(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("ip", ip.Properties.IP)
	d.Set("prefix", ip.Properties.Prefix)
	d.Set("location_uuid", ip.Properties.LocationUUID)
	d.Set("failover", ip.Properties.Failover)
	d.Set("status", ip.Properties.Status)
	d.Set("reverse_dns", ip.Properties.ReverseDNS)
	d.Set("location_country", ip.Properties.LocationCountry)
	d.Set("location_iata", ip.Properties.LocationIata)
	d.Set("location_name", ip.Properties.LocationName)
	d.Set("create_time", ip.Properties.CreateTime.String())
	d.Set("change_time", ip.Properties.ChangeTime.String())
	d.Set("delete_block", ip.Properties.DeleteBlock)
	d.Set("usage_in_minutes", ip.Properties.UsagesInMinutes)
	d.Set("current_price", ip.Properties.CurrentPrice)

	if err = d.Set("labels", ip.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleIpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.IPUpdateRequest{
		Name:       d.Get("name").(string),
		Failover:   d.Get("failover").(bool),
		ReverseDNS: d.Get("reverse_dns").(string),
		Labels:     convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	err := client.UpdateIP(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleIpRead(d, meta)
}

func resourceGridscaleIpv4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.IPCreateRequest{
		Family:       gsclient.IPv4Type,
		LocationUUID: d.Get("location_uuid").(string),
		Name:         d.Get("name").(string),
		Failover:     d.Get("failover").(bool),
		ReverseDNS:   d.Get("reverse_dns").(string),
		Labels:       convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	response, err := client.CreateIP(emptyCtx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new Ipv%v has been set to %v", requestBody.Family, response.ObjectUUID)

	return resourceGridscaleIpRead(d, meta)
}

func resourceGridscaleIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	ip, err := client.GetIP(emptyCtx, d.Id())
	if err != nil {
		return err
	}
	//Stop the server relating to this IP address if there is one
	//ip server relation is 1-1 relation
	if len(ip.Properties.Relations.Servers) == 1 {
		server := ip.Properties.Relations.Servers[0]
		deleteIPAction := func(ctx context.Context) error {
			err = client.DeleteIP(ctx, d.Id())
			return err
		}
		//DeleteIP requires the server to be off
		err = serverPowerStateList.runActionRequireServerOff(emptyCtx, client, server.ServerUUID, deleteIPAction)
		if err != nil {
			return err
		}
		return nil
	}
	return client.DeleteIP(emptyCtx, d.Id())
}
