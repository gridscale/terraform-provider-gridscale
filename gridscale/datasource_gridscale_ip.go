package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gridscale/gsclient-go"
)

func dataSourceGridscaleIp() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceGridscaleIpRead,
		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGridscaleIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	ip, err := client.GetIP(id)
	if err == nil {
		d.SetId(ip.Properties.ObjectUUID)
	}
	return err
}
