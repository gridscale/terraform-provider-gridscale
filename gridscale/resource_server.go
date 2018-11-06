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
		Update: resourceGridscaleServerUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters",
				Required:    true,
			},
			"memory": {
				Type:         schema.TypeInt,
				Description:  "The amount of server memory in GB.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cores": {
				Type:         schema.TypeInt,
				Description:  "The number of server cores.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"hardware_profile": {
				Type:        schema.TypeString,
				Description: "The number of server cores.",
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
			},
			"storage": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"power": {
				Type:        schema.TypeBool,
				Description: "The number of server cores.",
				Optional:    true,
				Default:     false,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
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
		Name:            d.Get("name").(string),
		Cores:           d.Get("cores").(int),
		Memory:          d.Get("memory").(int),
		LocationUuid:    d.Get("location_uuid").(string),
		HardwareProfile: d.Get("hardware_profile").(string),
	}

	createRequest.Relations.IsoImages = []interface{}{}
	createRequest.Relations.Networks = []gsclient.ServerNetwork{}
	createRequest.Relations.PublicIps = []gsclient.ServerIp{}

	createRequest.Relations.Storages = []gsclient.ServerStorage{}
	if attr, ok := d.GetOk("storage"); ok {
		for	_, value := range attr.([]interface{}) {
			storage := gsclient.ServerStorage{
				StorageUuid: value.(string),
			}
			createRequest.Relations.Storages = append(createRequest.Relations.Storages, storage)
		}
	}

	createRequest.Relations.PublicIps= []gsclient.ServerIp{}
	if attr, ok := d.GetOk("ip"); ok {
			ip := gsclient.ServerIp{
				IpaddrUuid: attr.(string),
			}
			createRequest.Relations.PublicIps = append(createRequest.Relations.PublicIps, ip)
	}

	createRequest.Relations.Networks = []gsclient.ServerNetwork{}
	if attr, ok := d.GetOk("network"); ok {
			network := gsclient.ServerNetwork{
				NetworkUuid: attr.(string),
				BootDevice:  true,
			}
			createRequest.Relations.Networks = []gsclient.ServerNetwork{network}
	}

	response, err := client.CreateServer(createRequest)

	if err != nil {
		return fmt.Errorf(
			"Error waiting for server (%s) to be created: %s", createRequest.Name, err)
	}

	d.SetId(response.ServerUuid)

	log.Printf("[DEBUG] The id for %s has been set to: %v", createRequest.Name, response.ServerUuid)

	power := d.Get("power").(bool)
	if power {
		client.StartServer(d.Id())
	}

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
	client := meta.(*gsclient.Client)
	requestBody := make(map[string]interface{})
	id := d.Id()

	if d.HasChange("name") {
		_, change := d.GetChange("name")
		requestBody["name"] = change.(string)
	}

	err := client.UpdateServer(id, requestBody)
	if err != nil {
		return err
	}

	if d.HasChange("power") {
		_, change := d.GetChange("power")
		power := change.(bool)
		if power {
			client.StartServer(id)
		}
	}

	return resourceGridscaleServerRead(d, meta)

}
