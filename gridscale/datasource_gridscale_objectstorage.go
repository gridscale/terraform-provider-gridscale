package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleObjectStorage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleObjectStorageRead,
		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
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

func dataSourceGridscaleObjectStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read object storage (%s) datasource-", id)

	objectStorage, err := client.GetObjectStorageAccessKey(context.Background(), id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	d.SetId(objectStorage.Properties.AccessKey)
	if err = d.Set("access_key", objectStorage.Properties.AccessKey); err != nil {
		return fmt.Errorf("%s error setting access_key: %v", errorPrefix, err)
	}
	if err = d.Set("secret_key", objectStorage.Properties.SecretKey); err != nil {
		return fmt.Errorf("%s error setting access_key: %v", errorPrefix, err)
	}

	return nil
}
