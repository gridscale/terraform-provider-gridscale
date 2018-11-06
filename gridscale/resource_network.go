package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceGridscaleNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleNetworkCreate,
		Read:   resourceGridscaleNetworkRead,
		Delete: resourceGridscaleNetworkDelete,
		Update: resourceGridscaleNetworkUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"l2security": {
				Type:        schema.TypeBool,
				Description: "MAC spoofing protection - filters layer2 and ARP traffic based on source MAC",
				Optional:    true,
				Default:     false,
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
				Description: "Path to the directory where the templated files will be written",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
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
			"public_net": {
				Type:        schema.TypeBool,
				Description: "Is the network public or not",
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
		},
	}
}

func resourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	network, err := client.GetNetwork(d.Id())

	d.Set("name", network.Properties.Name)
	d.Set("location_uuid", network.Properties.LocationUuid)
	d.Set("l2security", network.Properties.L2Security)
	d.Set("status", network.Properties.Status)
	d.Set("network_type", network.Properties.NetworkType)
	d.Set("location_country", network.Properties.LocationCountry)
	d.Set("location_iata", network.Properties.LocationIata)
	d.Set("location_name", network.Properties.LocationName)
	d.Set("public_net", network.Properties.PublicNet)
	d.Set("public_net", network.Properties.DeleteBlock)
	d.Set("create_time", network.Properties.CreateTime)
	d.Set("change_time", network.Properties.ChangeTime)

	log.Printf("Read the following: %v", network)
	return err
}

func resourceGridscaleNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := make(map[string]interface{})
	id := d.Id()

	if d.HasChange("name") {
		_, change := d.GetChange("name")
		requestBody["name"] = change.(string)
	}
	if d.HasChange("l2security") {
		_, change := d.GetChange("l2security")
		requestBody["l2security"] = change.(bool)
	}

	err := client.UpdateNetwork(id, requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["location_uuid"] = d.Get("location_uuid").(string)
	body["l2security"] = d.Get("l2security").(bool)

	response, err := client.CreateNetwork(body)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUuid)

	log.Printf("The id for network %v has been set to %v", d.Get("name").(string), response.ObjectUuid)

	return resourceGridscaleNetworkRead(d, meta)
}

func resourceGridscaleNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteNetwork(d.Id())

	return err
}
