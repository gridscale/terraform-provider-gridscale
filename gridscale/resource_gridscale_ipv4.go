package gridscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gridscale/gsclient-go/v2"
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
			Create: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
			Update: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
			Delete: schema.DefaultTimeout(time.Duration(GSCTimeoutSecs) * time.Second),
		},
	}
}

func resourceGridscaleIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read IP (%s) resource -", d.Id())
	ip, err := client.GetIP(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	if err = d.Set("ip", ip.Properties.IP); err != nil {
		return fmt.Errorf("%s error setting ip: %v", errorPrefix, err)
	}
	if err = d.Set("prefix", ip.Properties.Prefix); err != nil {
		return fmt.Errorf("%s error setting prefix: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", ip.Properties.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("failover", ip.Properties.Failover); err != nil {
		return fmt.Errorf("%s error setting failover: %v", errorPrefix, err)
	}
	if err = d.Set("status", ip.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("reverse_dns", ip.Properties.ReverseDNS); err != nil {
		return fmt.Errorf("%s error setting reverse_dns: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", ip.Properties.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", ip.Properties.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", ip.Properties.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", ip.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", ip.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("delete_block", ip.Properties.DeleteBlock); err != nil {
		return fmt.Errorf("%s error setting delete_block: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes", ip.Properties.UsagesInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", ip.Properties.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
	}

	if err = d.Set("labels", ip.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	return nil
}

func resourceGridscaleIpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update IP (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.IPUpdateRequest{
		Name:       d.Get("name").(string),
		Failover:   d.Get("failover").(bool),
		ReverseDNS: d.Get("reverse_dns").(string),
		Labels:     &labels,
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateIP(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleIpRead(d, meta)
}

func resourceGridscaleIpv4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.IPCreateRequest{
		Family:     gsclient.IPv4Type,
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

func resourceGridscaleIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete IP (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()

	ip, err := client.GetIP(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	//Stop the server relating to this IP address if there is one
	//ip server relation is 1-1 relation
	if len(ip.Properties.Relations.Servers) == 1 {
		server := ip.Properties.Relations.Servers[0]
		unlinkIPAction := func(ctx context.Context) error {
			err := client.UnlinkIP(ctx, server.ServerUUID, d.Id())
			return err
		}
		//DeleteIP requires the server to be off
		err = globalServerStatusList.runActionRequireServerOff(ctx, client, server.ServerUUID, false, unlinkIPAction)
		if err != nil {
			return fmt.Errorf("%s error: %v", errorPrefix, err)
		}
	}

	err = client.DeleteIP(ctx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
