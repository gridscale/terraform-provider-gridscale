package gridscale

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/nvthongswansea/gsclient-go"
)

func dataSourceGridscaleStorage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleStorageRead,
		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Description: "The capacity of a storage in GB.",
				Computed:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Helps to identify which datacenter an object belongs to.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "(one of storage, storage_high, storage_insane)",
				Computed:    true,
			},
			"license_product_no": {
				Type:        schema.TypeInt,
				Description: "If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).",
				Computed:    true,
			},
			"last_used_template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "status indicates the status of the object.",
				Computed:    true,
			},
			"create_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time the object was initially created.",
				Computed:    true,
			},
			"change_time": {
				Type:        schema.TypeString,
				Description: "Defines the date and time of the last object change.",
				Computed:    true,
			},
			"location_name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.",
				Computed:    true,
			},
			"location_country": {
				Type:        schema.TypeString,
				Description: "Formatted by the 2 digit country code (ISO 3166-2) of the host country.",
				Computed:    true,
			},
			"location_iata": {
				Type:        schema.TypeString,
				Description: "Uses IATA airport code, which works as a location identifier.",
				Computed:    true,
			},
			"current_price": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGridscaleStorageRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	storage, err := client.GetStorage(emptyCtx, id)

	if err == nil {
		d.SetId(storage.Properties.ObjectUUID)
		d.Set("change_time", storage.Properties.ChangeTime)
		d.Set("location_iata", storage.Properties.LocationIata)
		d.Set("status", storage.Properties.Status)
		d.Set("license_product_no", storage.Properties.LicenseProductNo)
		d.Set("location_country", storage.Properties.LocationCountry)
		d.Set("usage_in_minutes", storage.Properties.UsageInMinutes)
		d.Set("last_used_template", storage.Properties.LastUsedTemplate)
		d.Set("current_price", storage.Properties.CurrentPrice)
		d.Set("capacity", storage.Properties.Capacity)
		d.Set("location_uuid", storage.Properties.LocationUUID)
		d.Set("storage_type", storage.Properties.StorageType)
		d.Set("parent_uuid", storage.Properties.ParentUUID)
		d.Set("name", storage.Properties.Name)
		d.Set("location_name", storage.Properties.LocationName)
		d.Set("create_time", storage.Properties.CreateTime)

		if err = d.Set("labels", storage.Properties.Labels); err != nil {
			return fmt.Errorf("Error setting labels: %v", err)
		}
	}

	return err
}
