package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleSnapshotRead,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"storage_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Uuid of the storage used to create this snapshot",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the object",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status indicates the status of the object",
			},
			"location_country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the location",
			},
			"location_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the location",
			},
			"location_iata": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Uses IATA airport code, which works as a location identifier",
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Helps to identify which datacenter an object belongs to",
			},
			"usage_in_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total minutes the object has been running",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the date and time the object was initially created",
			},
			"change_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Defines the date and time of the last object change",
			},
			"license_product_no": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: `If a template has been used that requires a license key (e.g. Windows Servers) this shows 
the product_no of the license (see the /prices endpoint for more details)`,
			},
			"current_price": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The price for the current period since the last bill",
			},
			"capacity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The capacity of a storage/ISO Image/template/snapshot in GB",
			},
			"labels": {
				Type:        schema.TypeList,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscaleSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	storageUuid := d.Get("storage_uuid").(string)
	id := d.Get("resource_id").(string)
	errorPrefix := fmt.Sprintf("read snapshot (%s) datasource of storage (%s)-", id, storageUuid)

	snapshot, err := client.GetStorageSnapshot(context.Background(), storageUuid, id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	props := snapshot.Properties
	d.SetId(props.ObjectUUID)
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("location_country", props.LocationCountry); err != nil {
		return fmt.Errorf("%s error setting location_country: %v", errorPrefix, err)
	}
	if err = d.Set("location_name", props.LocationName); err != nil {
		return fmt.Errorf("%s error setting location_name: %v", errorPrefix, err)
	}
	if err = d.Set("location_iata", props.LocationIata); err != nil {
		return fmt.Errorf("%s error setting location_iata: %v", errorPrefix, err)
	}
	if err = d.Set("location_uuid", props.LocationUUID); err != nil {
		return fmt.Errorf("%s error setting location_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("usage_in_minutes", props.UsageInMinutes); err != nil {
		return fmt.Errorf("%s error setting usage_in_minutes: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}
	if err = d.Set("license_product_no", props.LicenseProductNo); err != nil {
		return fmt.Errorf("%s error setting license_product_no: %v", errorPrefix, err)
	}
	if err = d.Set("current_price", props.CurrentPrice); err != nil {
		return fmt.Errorf("%s error setting current_price: %v", errorPrefix, err)
	}
	if err = d.Set("capacity", props.Capacity); err != nil {
		return fmt.Errorf("%s error setting capacity: %v", errorPrefix, err)
	}
	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}

	return nil
}
