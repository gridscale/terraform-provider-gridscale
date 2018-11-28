package gridscale

import (
	"bitbucket.org/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func dataSourceGridscaleIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleIpRead,

		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGridscaleIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	ip, err := client.GetIp(id)

	if err == nil {
		d.SetId(ip.Properties.ObjectUuid)
		log.Printf("Found IP: %v", ip)
	}

	return err
}
