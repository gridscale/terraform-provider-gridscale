package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go"
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
				Description: "The capacity of a storage/ISO-Image/template/snapshot in GB",
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

	snapshot, err := client.GetStorageSnapshot(emptyCtx, storageUuid, id)

	if err == nil {
		props := snapshot.Properties
		d.SetId(props.ObjectUUID)
		d.Set("name", props.Name)
		d.Set("status", props.Status)
		d.Set("location_country", props.LocationCountry)
		d.Set("location_name", props.LocationName)
		d.Set("location_iata", props.LocationIata)
		d.Set("location_uuid", props.LocationUUID)
		d.Set("usage_in_minutes", props.UsageInMinutes)
		d.Set("create_time", props.CreateTime.String())
		d.Set("change_time", props.ChangeTime.String())
		d.Set("license_product_no", props.LicenseProductNo)
		d.Set("current_price", props.CurrentPrice)
		d.Set("capacity", props.Capacity)
		//Set labels
		if err = d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("Error setting labels: %v", err)
		}
	}

	return err
}
