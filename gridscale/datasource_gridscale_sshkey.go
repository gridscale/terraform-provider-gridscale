package gridscale

import (
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gridscale/gsclient-go"
)

func dataSourceGridscaleSshkey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleSshkeyRead,

		Schema: map[string]*schema.Schema{},
	}
}

func dataSourceGridscaleSshkeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	id := d.Get("resource_id").(string)
	sshkey, err := client.GetSshkey(id)
	if err == nil {
		d.SetId(sshkey.Properties.ObjectUUID)
	}
	return err
}
