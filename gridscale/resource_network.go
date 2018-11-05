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
				Description: "Name of the server",
				Required:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Path to the directory where the templated files will be written",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"l2security": {
				Type:        schema.TypeBool,
				Description: "Protects a network from MAC- and ARP-spoofing",
				Optional:    true,
				Default:     false,
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

	return client.UpdateNetwork(id, requestBody)
}

func resourceGridscaleNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["location_uuid"] = d.Get("location_uuid").(string)
	body["l2security"] = d.Get("l2security").(string)

	response, err := client.CreateNetwork(body)

	d.SetId(response.ObjectUuid)

	log.Printf("The id for network %v has been set to %v", d.Get("name").(string), response.ObjectUuid)

	return err
}

func resourceGridscaleNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteNetwork(d.Id())

	return err
}
