package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
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
		d.SetId(sshkey.Properties.ObjectUuid)
		log.Printf("Found SSH key: %v", sshkey)
	}

	return err
}
