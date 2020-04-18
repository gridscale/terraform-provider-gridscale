package gridscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gridscale/gsclient-go/v2"
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
				Description: "The capacity of a storage/ISO Image/template/snapshot in GB",
			},
			"storage_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Uuid of the storage used to create this snapshot",
			},
			"rollback": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Returns a storage to the state of the selected Snapshot.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the rollback request. Each rollback request has to have a unique id. ID can be any string value.",
						},
						"rollback_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(0 * time.Second),
			Update: schema.DefaultTimeout(0 * time.Second),
			Delete: schema.DefaultTimeout(0 * time.Second),
		},
	}
}

func resourceGridscaleSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	storageUuid := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("read snapshot (%s) resource of storage (%s)-", d.Id(), storageUuid)
	client := meta.(*gsclient.Client)
	snapshot, err := client.GetStorageSnapshot(context.Background(), storageUuid, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	props := snapshot.Properties
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

func resourceGridscaleSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	requestBody := gsclient.StorageSnapshotCreateRequest{
		Name:   d.Get("name").(string),
		Labels: convSOStrings(d.Get("labels").(*schema.Set).List()),
	}
	//set context with timeout when timeout is set
	ctx := context.Background()
	if d.Timeout(schema.TimeoutCreate) > zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
		defer cancel()
	}
	response, err := client.CreateStorageSnapshot(ctx, storageUUID, requestBody)
	if err != nil {
		return err
	}
	errorPrefix := fmt.Sprintf("rollback storage (%s) snapshot (%s) -", storageUUID, response.ObjectUUID)

	//Start rolling back if there are initially requests to rollback
	if attr, ok := d.GetOk("rollback"); ok {
		requests := make([]interface{}, 0)
		for _, requestProps := range attr.(*schema.Set).List() {
			rollbackReq := requestProps.(map[string]interface{})
			//Rollback
			log.Printf("Start rolling back storage %s with snapshot %s", storageUUID, response.ObjectUUID)
			err = client.RollbackStorage(
				ctx,
				storageUUID,
				response.ObjectUUID,
				gsclient.StorageRollbackRequest{
					Rollback: true,
				},
			)
			//Set status
			if err != nil {
				rollbackReq["status"] = err.Error()
			} else {
				rollbackReq["status"] = "success"
				log.Printf("Rolling back storage %s with snapshot %s SUCCESSFULLY", storageUUID, response.ObjectUUID)
			}
			//Set time of rollback request
			rollbackReq["rollback_time"] = time.Now().Format(timeLayout)
			requests = append(requests, rollbackReq)
		}
		//Apply value back to schema
		if err = d.Set("rollback", requests); err != nil {
			return fmt.Errorf("%s error setting rollback: %v", errorPrefix, err)
		}
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for snapshot %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleSnapshotRead(d, meta)
}

func resourceGridscaleSnapshotUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("update snapshot (%s) resource of storage (%s) -", d.Id(), storageUUID)

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.StorageSnapshotUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: &labels,
	}
	//set context with timeout when timeout is set
	ctx := context.Background()
	if d.Timeout(schema.TimeoutUpdate) > zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
		defer cancel()
	}
	err := client.UpdateStorageSnapshot(ctx, storageUUID, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	//Start rolling back if there are new requests to rollback
	if attr, ok := d.GetOk("rollback"); ok {
		requests := make([]interface{}, 0)
		for _, requestProps := range attr.(*schema.Set).List() {
			rollbackReq := requestProps.(map[string]interface{})
			//If `time` field of a request is set, it is the old request
			if rollbackReq["rollback_time"] == "" {
				//Rollback
				log.Printf("Start rolling back storage %s with snapshot %s", storageUUID, d.Id())
				err = client.RollbackStorage(
					ctx,
					storageUUID,
					d.Id(),
					gsclient.StorageRollbackRequest{
						Rollback: true,
					},
				)
				//Set status
				if err != nil {
					rollbackReq["status"] = err.Error()
				} else {
					rollbackReq["status"] = "success"
					log.Printf("Rolling back storage %s with snapshot %s SUCCESSFULLY", storageUUID, d.Id())
				}
				//Set time of rollback request
				rollbackReq["rollback_time"] = time.Now().Format(timeLayout)
			}
			requests = append(requests, rollbackReq)
		}
		//Apply value back to schema
		if err = d.Set("rollback", requests); err != nil {
			return fmt.Errorf("%s error setting rollback: %v", errorPrefix, err)
		}
	}
	return resourceGridscaleSnapshotRead(d, meta)
}

func resourceGridscaleSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("delete snapshot (%s) resource of storage (%s) -", d.Id(), storageUUID)
	//set context with timeout when timeout is set
	ctx := context.Background()
	if d.Timeout(schema.TimeoutDelete) > zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
		defer cancel()
	}
	err := client.DeleteStorageSnapshot(ctx, storageUUID, d.Id())
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
