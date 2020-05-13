package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleSshkey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleSshkeyRead,

		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true},
			"sshkey": {
				Type:        schema.TypeString,
				Description: "sshkey_string is the OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1)",
				Computed:    true},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "The date and time the object was initially created",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "The date and time of the last object change",
				Computed:    true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscaleSshkeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read SSH key (%s) datasource -", id)

	sshkey, err := client.GetSshkey(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(sshkey.Properties.ObjectUUID)
	if err = d.Set("name", sshkey.Properties.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("sshkey", sshkey.Properties.Sshkey); err != nil {
		return fmt.Errorf("%s error setting sshkey: %v", errorPrefix, err)
	}
	if err = d.Set("status", sshkey.Properties.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", sshkey.Properties.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", sshkey.Properties.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("labels", sshkey.Properties.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	return nil
}
