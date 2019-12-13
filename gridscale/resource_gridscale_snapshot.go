package gridscale

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"

	"github.com/gridscale/gsclient-go"
)

func resourceGridscaleStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGridscaleSnapshotRead,
		Create: resourceGridscaleSnapshotCreate,
		Update: resourceGridscaleSnapshotUpdate,
		Delete: resourceGridscaleSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
			"storage_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Uuid of the storage used to create this snapshot",
			},
			"rollback": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Returns a storage to the state of the selected Snapshot.",
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGridscaleSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	storageUuid := d.Get("storage_uuid").(string)
	client := meta.(*gsclient.Client)
	snapshot, err := client.GetStorageSnapshot(emptyCtx, storageUuid, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := snapshot.Properties
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
	return nil
}

func resourceGridscaleSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	requestBody := gsclient.StorageSnapshotCreateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	response, err := client.CreateStorageSnapshot(emptyCtx, storageUUID, requestBody)
	if err != nil {
		return err
	}
	//Check if rollback is requested
	rollback := d.Get("rollback").(bool)
	if rollback {
		log.Printf("Rolling back storage %s using snapshot %s", storageUUID, response.ObjectUUID)
		err = client.RollbackStorage(
			emptyCtx,
			storageUUID,
			response.ObjectUUID,
			gsclient.StorageRollbackRequest{
				Rollback: rollback,
			},
		)
		if err != nil {
			return err
		}
		log.Printf("Rolling back storage %s using snapshot %s SUCCESSFULLY", storageUUID, response.ObjectUUID)
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for snapshot %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleSnapshotRead(d, meta)
}

func resourceGridscaleSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.StorageSnapshotUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	storageUUID := d.Get("storage_uuid").(string)
	err := client.UpdateStorageSnapshot(emptyCtx, storageUUID, d.Id(), requestBody)
	if err != nil {
		return err
	}
	//Check if rollback is requested
	rollback := d.Get("rollback").(bool)
	rollbackHasChanged := d.HasChange("rollback")
	if rollback && rollbackHasChanged {
		log.Printf("Rolling back storage %s using snapshot %s", storageUUID, d.Id())
		err = client.RollbackStorage(
			emptyCtx,
			storageUUID,
			d.Id(),
			gsclient.StorageRollbackRequest{
				Rollback: rollback,
			},
		)
		if err != nil {
			return err
		}
		log.Printf("Rolling back storage %s using snapshot %s SUCCESSFULLY", storageUUID, d.Id())
	}
	return resourceGridscaleSnapshotRead(d, meta)
}

func resourceGridscaleSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	return client.DeleteStorageSnapshot(emptyCtx, d.Get("storage_uuid").(string), d.Id())
}
