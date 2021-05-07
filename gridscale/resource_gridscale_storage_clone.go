package gridscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gridscale/gsclient-go/v3"
)

func resourceGridscaleStorageClone() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleStorageCloneCreate,
		Read:   resourceGridscaleStorageRead,
		Delete: resourceGridscaleStorageDelete,
		Update: resourceGridscaleStorageUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"source_storage_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "ID of the storage instance that will be cloned.",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.",
				Optional:    true,
				Computed:    true,
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

func resourceGridscaleStorageCloneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	srcStorageID := d.Get("source_storage_id").(string)
	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CloneStorage(ctx, srcStorageID)
	if err != nil {
		return err
	}

	d.SetId(response.ObjectUUID)

	log.Printf("A new storage %s has been created by cloning storage %s", response.ObjectUUID, srcStorageID)

	// If the user wants a new name and new capacity for the storage clone
	// instead of the default values inherited from the source storage,
	// change name and capacity of the storage clone to the desired ones.
	sourceStorage, err := client.GetStorage(ctx, srcStorageID)
	if err != nil {
		return err
	}
	if sourceStorage.Properties.Name != d.Get("name").(string) ||
		sourceStorage.Properties.Capacity != d.Get("capacity").(int) ||
		sourceStorage.Properties.StorageType != d.Get("storage_type").(string) {
		err = resourceGridscaleStorageUpdate(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceGridscaleStorageRead(d, meta)
}
