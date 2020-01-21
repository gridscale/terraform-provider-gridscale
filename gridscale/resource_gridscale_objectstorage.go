package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
)

func resourceGridscaleObjectStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleObjectStorageCreate,
		Read:   resourceGridscaleObjectStorageRead,
		Delete: resourceGridscaleObjectStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Description: "The object storage secret_key.",
				Computed:    true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The object storage access_key.",
				Computed:    true,
			},
		},
	}
}

func resourceGridscaleObjectStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read object storage (%s) resource -", d.Id())
	objectStorage, err := client.GetObjectStorageAccessKey(emptyCtx, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	if err = d.Set("access_key", objectStorage.Properties.AccessKey); err != nil {
		return fmt.Errorf("%s error setting access_key: %v", errorPrefix, err)
	}
	if err = d.Set("secret_key", objectStorage.Properties.SecretKey); err != nil {
		return fmt.Errorf("%s error setting secret_key: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscaleObjectStorageCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	response, err := client.CreateObjectStorageAccessKey(emptyCtx)
	if err != nil {
		return err
	}

	d.SetId(response.AccessKey.AccessKey)

	log.Printf("The id for the new object storage has been set to %v", response.AccessKey.AccessKey)
	return resourceGridscaleObjectStorageRead(d, meta)
}

func resourceGridscaleObjectStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete object storage (%s) resource -", d.Id())
	err := client.DeleteObjectStorageAccessKey(emptyCtx, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
