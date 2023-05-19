package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGridscaleObjectStorage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleObjectStorageRead,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"comment": {
				Type:        schema.TypeString,
				Description: "Comment for the access_key.",
				Computed:    true,
			},
			"user_uuid": {
				Type: schema.TypeString,
				Description: `If a user_uuid is set, a user-specific key will get created. 
				If no user_uuid is set along a user with write-access to the contract will still only create 
				a user-specific key for themselves while a user with admin-access to the contract will create 
				a contract-level admin key.`,
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
	}
}

func dataSourceGridscaleObjectStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read object storage (%s) datasource-", id)

	objectStorage, err := client.GetObjectStorageAccessKey(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(objectStorage.Properties.AccessKey)
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
		return fmt.Errorf("%s error setting access_key: %v", errorPrefix, err)
	}

	return nil
}
