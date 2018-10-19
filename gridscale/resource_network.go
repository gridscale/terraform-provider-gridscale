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
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        	schema.TypeString,
				Description: 	"Name of the server",
				Required:    	true,
				ForceNew:	 	true,
			},
			"location_uuid": {
				Type:        	schema.TypeString,
				Description: 	"Path to the directory where the templated files will be written",
				Optional:    	true,
				ForceNew:		true,
				Default:	 	"45ed677b-3702-4b36-be2a-a2eab9827950",
			},
		},
	}
}

func resourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	network, err := client.ReadNetwork(d.Id())

	d.Set("name", network.Body.Name)
	d.Set("location_uuid", network.Body.LocationUuid)


	log.Printf("Read the following: %v", network)
	return err
}

func resourceGridscaleNetworkUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGridscaleNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["location_uuid"] = d.Get("location_uuid").(string)

	response, err := client.CreateNetwork(body)

	d.SetId(response.ObjectUuid)

	log.Printf("The id for blah has been set to %v", response)

	return err
}

func resourceGridscaleNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DestroyNetwork(d.Id())

	return err
}