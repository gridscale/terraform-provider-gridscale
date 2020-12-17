package gridscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go/v3"
)

func resourceGridscaleStorageImport() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleStorageImportFromBackup,
		Read:   resourceGridscaleStorageRead,
		Delete: resourceGridscaleStorageDelete,
		Update: resourceGridscaleStorageUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"storage_backup_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "ID of the storage backup that will be used to create a new storage from.",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Required:    true,
			},
			"capacity": {
				Type:        schema.TypeInt,
				Description: "The capacity of a storage in GB.",
				Optional:    true,
				Computed:    true,
			},
			"location_uuid": {
				Type:        schema.TypeString,
				Description: "Identifies the data center this object belongs to.",
				Computed:    true,
			},
			"storage_type": {
				Type:        schema.TypeString,
				Description: "(one of storage, storage_high, storage_insane)",
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					valid := false
					for _, stype := range storageTypes {
						if v.(string) == stype {
							valid = true
							break
						}
					}
					if !valid {
						errors = append(errors, fmt.Errorf("%v is not a valid storage type. Valid types are: %v", v.(string), strings.Join(storageTypes, ",")))
					}
					return
				},
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
				Description: "The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.",
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
			"usage_in_minutes": {
				Type:     schema.TypeInt,
				Computed: true,
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

func resourceGridscaleStorageImportFromBackup(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageBackupID := d.Get("storage_backup_id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateStorageFromBackup(ctx, storageBackupID, d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("A new storage %s has been created from backup %s", response.ObjectUUID, storageBackupID)

	// If the user wants a new storage type and new capacity for the storage clone
	// instead of the default values inherited from the storage backup,
	// change storage type and capacity of the storage clone to the desired ones.
	newStorage, err := client.GetStorage(ctx, response.ObjectUUID)
	if err != nil {
		return err
	}
	if newStorage.Properties.Capacity != d.Get("capacity").(int) ||
		newStorage.Properties.StorageType != d.Get("storage_type").(string) {
		err = resourceGridscaleStorageUpdate(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceGridscaleStorageRead(d, meta)
}
