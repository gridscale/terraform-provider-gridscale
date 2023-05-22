package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"
)

func resourceGridscaleObjectStorage() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleObjectStorageCreate,
		Update: resourceGridscaleObjectStorageUpdate,
		Read:   resourceGridscaleObjectStorageRead,
		Delete: resourceGridscaleObjectStorageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the access_key.",
				Optional:    true,
				Computed:    true,
			},
			"user_uuid": {
				Type: schema.TypeString,
				Description: `If a user_uuid is set, a user-specific key will get created. 
				If no user_uuid is set along a user with write-access to the contract will still only create 
				a user-specific key for themselves while a user with admin-access to the contract will create 
				a contract-level admin key.`,
				Optional: true,
				Computed: true,
			},
			"access_key": {
				Type:        schema.TypeString,
				Description: "The object storage secret_key.",
				Computed:    true,
				Sensitive:   true,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Description: "The object storage access_key.",
				Computed:    true,
				Sensitive:   true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleObjectStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("read object storage (%s) resource -", d.Id())
	objectStorage, err := client.GetObjectStorageAccessKey(context.Background(), d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	if err = d.Set("comment", objectStorage.Properties.Comment); err != nil {
		return fmt.Errorf("%s error setting comment: %v", errorPrefix, err)
	}
	if err = d.Set("user_uuid", objectStorage.Properties.UserUUID); err != nil {
		return fmt.Errorf("%s error setting user_uuid: %v", errorPrefix, err)
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

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	requestBody := gsclient.ObjectStorageAccessKeyCreateRequest{}
	if comment, ok := d.GetOk("comment"); ok {
		requestBody.Comment = comment.(string)
	}
	if userUUID, ok := d.GetOk("user_uuid"); ok {
		requestBody.UserUUID = userUUID.(string)
	}
	response, err := client.AdvancedCreateObjectStorageAccessKey(ctx, requestBody)
	if err != nil {
		return err
	}

	d.SetId(response.AccessKey.AccessKey)

	log.Printf("The id for the new object storage has been set to %v", response.AccessKey.AccessKey)
	return resourceGridscaleObjectStorageRead(d, meta)
}

func resourceGridscaleObjectStorageUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("update object storage (%s) resource -", d.Id())
	requestBody := gsclient.ObjectStorageAccessKeyUpdateRequest{}
	if d.HasChange("comment") {
		newComment, ok := d.Get("comment").(string)
		if ok {
			requestBody.Comment = &newComment
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateObjectStorageAccessKey(ctx, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	return resourceGridscaleSshkeyRead(d, meta)
}

func resourceGridscaleObjectStorageDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	errorPrefix := fmt.Sprintf("delete object storage (%s) resource -", d.Id())

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.SuppressHTTPErrorCodes(
		client.DeleteObjectStorageAccessKey(ctx, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
