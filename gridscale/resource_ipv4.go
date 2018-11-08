package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceGridscaleIpv4() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleIpv4Create,
		Read:   resourceGridscaleIpRead,
		Delete: resourceGridscaleIpDelete,
		Update: resourceGridscaleIpUpdate,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Description: "Defines the IP Address.",
				Computed:    true,
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
				Description: "The date and time the object was initially created",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change",
				Computed:    true,
			},
		},
	}
}

func resourceGridscaleIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	ip, err := client.GetIp(d.Id())

	d.Set("ip", ip.Properties.Ip)
	d.Set("prefix", ip.Properties.Prefix)
	d.Set("family", ip.Properties.Family)
	d.Set("location_uuid", ip.Properties.LocationUuid)
	d.Set("failover", ip.Properties.Failover)
	d.Set("status", ip.Properties.Status)
	d.Set("reverse_dns", ip.Properties.ReverseDns)
	d.Set("location_country", ip.Properties.LocationCountry)
	d.Set("location_iata", ip.Properties.LocationIata)
	d.Set("location_name", ip.Properties.LocationName)
	d.Set("create_time", ip.Properties.CreateTime)
	d.Set("change_time", ip.Properties.ChangeTime)

	log.Printf("Read the following: %v", ip)
	return err
}

func resourceGridscaleIpUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := make(map[string]interface{})
	id := d.Id()

	if d.HasChange("failover") {
		_, change := d.GetChange("failover")
		requestBody["failover"] = change.(bool)
	}
	if d.HasChange("reverse_dns") {
		_, change := d.GetChange("reverse_dns")
		requestBody["reverse_dns"] = change.(string)
	}

	err := client.UpdateIp(id, requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleIpRead(d, meta)
}

func resourceGridscaleIpv4Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["family"] = 4
	body["location_uuid"] = d.Get("location_uuid").(string)
	body["failover"] = d.Get("failover").(bool)
	reversedns := d.Get("reverse_dns").(string)
	if reversedns != "" {
		body["reverse_dns"] = reversedns
	}

	response, err := client.CreateIp(body)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUuid)

	log.Printf("The id for the new Ipv%v has been set to %v", body["family"], response.ObjectUuid)

	return resourceGridscaleIpRead(d, meta)
}

func resourceGridscaleIpDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteIp(d.Id())

	return err
}
