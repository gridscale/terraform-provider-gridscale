package gridscale

import (
	"../gsclient"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceGridscaleSshkey() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleSshkeyCreate,
		Read:   resourceGridscaleSshkeyRead,
		Delete: resourceGridscaleSshkeyDelete,
		Update: resourceGridscaleSshkeyUpdate,
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
		},
	}
}

func resourceGridscaleSshkeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	sshkey, err := client.GetSshkey(d.Id())
	if err != nil {
		if requestError, ok := err.(*gsclient.RequestError); ok {
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

	log.Printf("Read the following: %v", sshkey)
	return nil
}

func resourceGridscaleSshkeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := make(map[string]interface{})
	id := d.Id()

	if d.HasChange("name") {
		_, change := d.GetChange("name")
		requestBody["name"] = change.(string)
	}
	if d.HasChange("sshkey") {
		_, change := d.GetChange("sshkey")
		requestBody["sshkey"] = change.(string)
	}

	err := client.UpdateSshkey(id, requestBody)
	if err != nil {
		return err
	}

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleSshkeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	body := make(map[string]interface{})
	body["name"] = d.Get("name").(string)
	body["sshkey"] = d.Get("sshkey").(string)

	response, err := client.CreateSshkey(body)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUuid)

	log.Printf("The id for the new SSH Key has been set to %v", response.ObjectUuid)

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleSshkeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	err := client.DeleteSshkey(d.Id())

	return err
}
