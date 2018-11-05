package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceGridscaleNetwork() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleNetworkRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceGridscaleNetworkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	storage, err := client.GetNetwork(id)

	if err == nil {
		d.SetId(storage.Properties.ObjectUuid)
		log.Printf("Found storage: %v", storage)
	}

	return err
}
