package gridscale

import (
	"../gsclient"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceGridscaleServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleServerCreate,
		Read:   resourceGridscaleServerRead,
		Delete: resourceGridscaleServerDelete,
		Update:	resourceGridscaleServerUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        	schema.TypeString,
				Description: 	"Name of the server",
				Required:    	true,
			},
			"memory": {
				Type:         	schema.TypeInt,
				Description:  	"Memory in gigabytes",
				Required:    	true,
				ValidateFunc:	validation.NoZeroValues,
			},
			"cores": {
				Type:        	schema.TypeInt,
				Description: 	"Amount of CPU cores",
				Required:    	true,
				ValidateFunc:	validation.NoZeroValues,
			},
			"location_uuid": {
				Type:        	schema.TypeString,
				Description: 	"Path to the directory where the templated files will be written",
				Optional:    	true,
				ForceNew:		true,
				Default:	 	"45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"hardware_profile": {
				Type:        	schema.TypeString,
				Computed:		true,
			},
			"storage": {
				Type:			schema.TypeString,
				Optional:		true,
			},
			"network": {
				Type:			schema.TypeString,
				Optional:		true,
			},
			"power": {
				Type:			schema.TypeBool,
				Computed:		true,
			},
			"current_price": {
				Type:			schema.TypeFloat,
				Computed:		true,
			},
		},
	}
}

func resourceGridscaleServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	server, err := client.GetServer(d.Id())

	d.Set("name", server.Properties.Name)
	d.Set("memory", server.Properties.Memory)
	d.Set("cores", server.Properties.Cores)
	d.Set("hardware_profile", server.Properties.HardwareProfile)
	d.Set("location_uuid", server.Properties.LocationUuid)
	d.Set("power", server.Properties.Power)
	d.Set("current_price", server.Properties.CurrentPrice)


	log.Printf("Read the following: %v", server)

	return err
}

func resourceGridscaleServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	createRequest := gsclient.ServerCreateRequest{
		Name:			d.Get("name").(string),
		Cores:			d.Get("cores").(int),
		Memory:			d.Get("memory").(int),
		LocationUuid:	d.Get("location_uuid").(string),
	}

	createRequest.Relations.IsoImages = []interface{}{}
	createRequest.Relations.Networks = []interface{}{}
	createRequest.Relations.PublicIps = []interface{}{}

	if attr, ok := d.GetOk("storage"); ok {
		storage := gsclient.ServerStorage{
			StorageUuid:	attr.(string),
			BootDevice:		true,
		}
		createRequest.Relations.Storages = []gsclient.ServerStorage{storage}
	} else {
		createRequest.Relations.Storages = []gsclient.ServerStorage{}
	}

	if attr, ok := d.GetOk("network"); ok {
		network := gsclient.ServerStorage{
			StorageUuid:	attr.(string),
			BootDevice:		true,
		}
		createRequest.Relations.Networks = []interface{}{network}
	} else {
		createRequest.Relations.Networks = []interface{}{}
	}

	//if attr, ok := d.GetOk("storage"); ok {
	//	test := attr.(string)
	//}

	response, err := client.CreateServer(createRequest)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", createRequest.Name, err)
	}

	d.SetId(response.ServerUuid)

	log.Printf("[DEBUG] The id for %s has been set to: %v", createRequest.Name, response.ServerUuid)

	return resourceGridscaleServerRead(d, meta)
}

func resourceGridscaleServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Id()
	err := client.StopServer(id)
	if err != nil {
		return err
	}
	err = client.DeleteServer(id)

	return err
}

func resourceGridscaleServerUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}