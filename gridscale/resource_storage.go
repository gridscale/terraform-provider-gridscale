package gridscale

import (
	"../gsclient"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
)

func resourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleStorageCreate,
		Read:   resourceGridscaleStorageRead,
		Delete: resourceGridscaleStorageDelete,
		Update: resourceGridscaleStorageUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"capacity": {
				Type:         schema.TypeInt,
				Description:  "The capacity of a storage in GB.",
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Optional:    true,
				ForceNew:    true,
				Default:     "45ed677b-3702-4b36-be2a-a2eab9827950",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "(one of storage, storage_high, storage_insane)",
				Optional:    true,
				ForceNew:    true,
				Default:     "storage",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if v.(string) != "storage" && v.(string) != "storage_high" && v.(string) != "storage_insane" {
						errors = append(errors, fmt.Errorf("Storage type must be either storage, storage_high or storage_insane"))
					}
					return
				},
			},
			"license_product_no": {
				Type:        schema.TypeInt,
				Description: "If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time of the last object change.",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Formatted by the 2 digit country code (ISO 3166-2) of the host country.",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier.",
				Computed:    true,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Description: "The price for the current period since the last bill.",
				Computed:    true,
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
	d.Set("status", storage.Properties.Status)
	d.Set("create_time", storage.Properties.CreateTime)
	d.Set("change_time", storage.Properties.ChangeTime)
	d.Set("location_name", storage.Properties.LocationName)
	d.Set("location_country", storage.Properties.LocationCountry)
	d.Set("location_iata", storage.Properties.LocationIata)
	d.Set("current_price", storage.Properties.CurrentPrice)

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

	err := client.UpdateStorage(id, requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["capacity"] = d.Get("capacity").(int)
	body["location_uuid"] = d.Get("location_uuid").(string)
	body["storage_type"] = d.Get("storage_type").(string)

	response, err := client.CreateStorage(body)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUuid)

	log.Printf("The id for storage %s has been set to %v", body["name"], response)

	return resourceGridscaleStorageRead(d, meta)
}

func resourceGridscaleStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteStorage(d.Id())

	return err
}
