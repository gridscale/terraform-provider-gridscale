package gridscale

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/nvthongswansea/gsclient-go"
)

func resourceGridscaleSshkey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleSshkeyCreate,
		Read:   resourceGridscaleSshkeyRead,
		Delete: resourceGridscaleSshkeyDelete,
		Update: resourceGridscaleSshkeyUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Required:    true,
			},
			"sshkey": {
				Type:        schema.TypeString,
				Description: "sshkey_string is the OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1)",
				Required:    true,
			},
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
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGridscaleSshkeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	sshkey, err := client.GetSshkey(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}

	d.Set("name", sshkey.Properties.Name)
	d.Set("sshkey", sshkey.Properties.Sshkey)
	d.Set("status", sshkey.Properties.Status)
	d.Set("create_time", sshkey.Properties.CreateTime)
	d.Set("change_time", sshkey.Properties.ChangeTime)

	if err = d.Set("labels", sshkey.Properties.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}

	return nil
}

func resourceGridscaleSshkeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.SshkeyUpdateRequest{
		Name:   d.Get("name").(string),
		Sshkey: d.Get("sshkey").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	err := client.UpdateSshkey(emptyCtx, d.Id(), requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleSshkeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	requestBody := gsclient.SshkeyCreateRequest{
		Name:   d.Get("name").(string),
		Sshkey: d.Get("sshkey").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	response, err := client.CreateSshkey(emptyCtx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new SSH Key has been set to %v", response.ObjectUUID)

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleSshkeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteSshkey(emptyCtx, d.Id())

	return err
}
