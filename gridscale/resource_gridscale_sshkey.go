package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
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
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
			},
			"sshkey": {
				Type:        schema.TypeString,
				Description: "sshkey_string is the OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1)",
				Required:    true,
				StateFunc: func(val interface{}) string {
					return strings.TrimSpace((val.(string)))
				},
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleSshkeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read SSH key (%s) resource -", d.Id())
	sshkey, err := client.GetSshkey(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

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

func resourceGridscaleSshkeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update SSH key (%s) resource -", d.Id())

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	pubKey := d.Get("sshkey").(string)
	requestBody := gsclient.SshkeyUpdateRequest{
		Name:   d.Get("name").(string),
		Sshkey: strings.TrimSpace(pubKey),
		Labels: &labels,
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateSshkey(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleSshkeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	pubKey := d.Get("sshkey").(string)
	requestBody := gsclient.SshkeyCreateRequest{
		Name:   d.Get("name").(string),
		Sshkey: strings.TrimSpace(pubKey),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateSshkey(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("The id for the new SSH Key has been set to %v", response.ObjectUUID)

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleSshkeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete SSH key (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.SuppressHTTPErrorCodes(
		client.DeleteSshkey(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
