package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gridscale/gsclient-go"
)

func dataSourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceGridscaleStorageRead,
		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGridscaleStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	storage, err := client.GetStorage(id)
	if err == nil {
		d.SetId(storage.Properties.ObjectUUID)
	}
	return err
}
