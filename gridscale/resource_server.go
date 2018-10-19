package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceGridscaleServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleServerCreate,
		Read:   resourceGridscaleServerRead,
		Delete: resourceGridscaleServerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        	schema.TypeString,
				Description: 	"Name of the server",
				Required:    	true,
				ForceNew:	 	true,
			},
			"memory": {
				Type:         	schema.TypeInt,
				Description:  	"Memory in gigabytes",
				Required:    	true,
				ForceNew:	 	true,
				ValidateFunc:	validation.NoZeroValues,
			},
			"cores": {
				Type:        	schema.TypeInt,
				Description: 	"Amount of CPU cores",
				Required:    	true,
				ForceNew:	 	true,
				ValidateFunc:	validation.NoZeroValues,
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

func resourceGridscaleServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	server, err := client.ReadServer(d.Id())

	d.Set("name", server.Body.Name)
	d.Set("memory", server.Body.Memory)
	d.Set("cores", server.Body.Cores)
	d.Set("hardware_profile", server.Body.HardwareProfile)
	d.Set("location_uuid", server.Body.LocationUuid)


	log.Printf("Read the following: %v", server)
	return err
}

func resourceGridscaleServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["memory"] = d.Get("memory").(int)
	body["cores"] = d.Get("cores").(int)
	body["location_uuid"] = d.Get("location_uuid").(string)

	response, err := client.CreateServer(body)

	d.SetId(response.ObjectUuid)

	log.Printf("The id for blah has been set to %v", response)

	return err
}

func resourceGridscaleServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DestroyServer(d.Id())

	return err
}