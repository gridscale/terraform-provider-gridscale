package gridscale

import (
	"bitbucket.org/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/schema"
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
	}

	return err
}
