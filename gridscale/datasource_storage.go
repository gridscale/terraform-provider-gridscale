package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleStorageRead,

		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGridscaleStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	storage, err := client.GetStorage(id)

	if err == nil {
		d.SetId(storage.Properties.ObjectUuid)
		log.Printf("Found storage: %v", storage)
	}

	return err
}
