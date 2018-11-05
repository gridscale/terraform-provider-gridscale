package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleStorageCreate,
		Read:   resourceGridscaleStorageRead,
		Delete: resourceGridscaleStorageDelete,
		Update:	resourceGridscaleStorageUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the server",
				Required:    true,
			},
			"capacity": {
				Type:         schema.TypeInt,
				Description:  "Storage capacity in gigabytes",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Path to the directory where the templated files will be written",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
		},
	}
}

func resourceGridscaleStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storage, err := client.GetStorage(d.Id())

	d.Set("name", storage.Properties.Name)
	d.Set("capacity", storage.Properties.Capacity)
	d.Set("location_uuid", storage.Properties.LocationUuid)

	log.Printf("Read the following: %v", storage)
	return err
}

func resourceGridscaleStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := make(map[string]interface{})
	id := d.Id()

	if d.HasChange("name") {
		_, change := d.GetChange("name")
		requestBody["name"] = change.(string)
	}

	if d.HasChange("capacity") {
		_, change := d.GetChange("capacity")
		requestBody["capacity"] = change.(int)
	}

	return client.UpdateStorage(id, requestBody)
}

func resourceGridscaleStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["capacity"] = d.Get("capacity").(int)
	body["location_uuid"] = d.Get("location_uuid").(string)

	response, err := client.CreateStorage(body)

	d.SetId(response.ObjectUuid)

	log.Printf("The id for storage %s has been set to %v", body["name"], response)

	return err
}

func resourceGridscaleStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteStorage(d.Id())

	return err
}
