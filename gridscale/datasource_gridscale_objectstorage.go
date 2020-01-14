package gridscale

import (
	"fmt"
	"github.com/gridscale/gsclient-go"
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

	objectStorage, err := client.GetObjectStorageAccessKey(emptyCtx, id)

	if err == nil {
		d.SetId(objectStorage.Properties.AccessKey)
		if err = d.Set("access_key", objectStorage.Properties.AccessKey); err != nil {
			return fmt.Errorf("error setting access_key: %v", err)
		}
		if err = d.Set("secret_key", objectStorage.Properties.SecretKey); err != nil {
			return fmt.Errorf("error setting access_key: %v", err)
		}
	}

	return err
}
