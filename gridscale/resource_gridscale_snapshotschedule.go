package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
)

func resourceGridscaleStorageSnapshotSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleSnapshotScheduleCreate,
		Read:   resourceGridscaleSnapshotScheduleRead,
		Delete: resourceGridscaleSnapshotScheduleDelete,
		Update: resourceGridscaleSnapshotScheduleUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The human-readable name of the object",
			},
			"next_runtime": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The date and time that the snapshot schedule will be run",
			},
			"next_runtime_computed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the snapshot schedule will be run. This date and time is computed by gridscale's server.",
			},
			"keep_snapshots": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The amount of Snapshots to keep before overwriting the last created Snapshot",
			},
			"run_interval": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(60),
				Description:  "The interval at which the schedule will run (in minutes)",
			},
			"storage_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the storage used to create snapshots",
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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status indicates the status of the object",
			},
			"labels": {
				Type:        schema.TypeSet,
				Description: "List of labels.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"snapshot": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Related snashots",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleSnapshotScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("read snapshot schedule (%s) resource of storage (%s)-", d.Id(), storageUUID)
	scheduler, err := client.GetStorageSnapshotSchedule(context.Background(), storageUUID, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	props := scheduler.Properties
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("next_runtime_computed", props.NextRuntime.Format(timeLayout)); err != nil {
		return fmt.Errorf("%s error setting next_runtime_computed: %v", errorPrefix, err)
	}
	if err = d.Set("keep_snapshots", props.KeepSnapshots); err != nil {
		return fmt.Errorf("%s error setting keep_snapshots: %v", errorPrefix, err)
	}
	if err = d.Set("run_interval", props.RunInterval); err != nil {
		return fmt.Errorf("%s error setting run_interval: %v", errorPrefix, err)
	}
	if err = d.Set("storage_uuid", props.StorageUUID); err != nil {
		return fmt.Errorf("%s error setting storage_uuid: %v", errorPrefix, err)
	}
	if err = d.Set("create_time", props.CreateTime.String()); err != nil {
		return fmt.Errorf("%s error setting create_time: %v", errorPrefix, err)
	}
	if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
		return fmt.Errorf("%s error setting change_time: %v", errorPrefix, err)
	}

	//Get snapshots
	snapshots := make([]interface{}, 0)
	for _, value := range props.Relations.Snapshots {
		snapshots = append(snapshots, map[string]interface{}{
			"name":        value.Name,
			"object_uuid": value.ObjectUUID,
			"create_time": value.CreateTime.String(),
		})
	}
	if err = d.Set("snapshot", snapshots); err != nil {
		return fmt.Errorf("%s error setting snapshots: %v", errorPrefix, err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("%s error setting labels: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscaleSnapshotScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.StorageSnapshotScheduleCreateRequest{
		Name:          d.Get("name").(string),
		Labels:        convSOStrings(d.Get("labels").(*schema.Set).List()),
		RunInterval:   d.Get("run_interval").(int),
		KeepSnapshots: d.Get("keep_snapshots").(int),
	}
	if strings.TrimSpace(d.Get("next_runtime").(string)) != "" {
		nextRuntime, err := time.Parse(timeLayout, d.Get("next_runtime").(string))
		if err != nil {
			return err
		}
		requestBody.NextRuntime = &gsclient.GSTime{Time: nextRuntime}
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateStorageSnapshotSchedule(ctx, d.Get("storage_uuid").(string), requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for snapshot schedule %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleSnapshotScheduleRead(d, meta)
}

func resourceGridscaleSnapshotScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("update snapshot schedule (%s) resource of storage (%s)-", d.Id(), storageUUID)

	labels := convSOStrings(d.Get("labels").(*schema.Set).List())
	requestBody := gsclient.StorageSnapshotScheduleUpdateRequest{
		Name:          d.Get("name").(string),
		Labels:        &labels,
		RunInterval:   d.Get("run_interval").(int),
		KeepSnapshots: d.Get("keep_snapshots").(int),
	}
	// NOTE: remember to check if next_runtime is changed,
	// otherwise tf will force to update next_runtime every time Update is executed.
	// This is a bad behavior. Because if next_runtime_computed is different from
	// next_runtime (since it might be changed outside of tf), and next_runtime is not changed (by the user);
	// tf should not put the "old" next_runtime to the update request.
	if d.HasChange("next_runtime") {
		if strings.TrimSpace(d.Get("next_runtime").(string)) != "" {
			nextRuntime, err := time.Parse(timeLayout, d.Get("next_runtime").(string))
			if err != nil {
				return fmt.Errorf("%s error: %v", errorPrefix, err)
			}
			requestBody.NextRuntime = &gsclient.GSTime{Time: nextRuntime}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err := client.UpdateStorageSnapshotSchedule(ctx, storageUUID, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleSnapshotScheduleRead(d, meta)
}

func resourceGridscaleSnapshotScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("delete snapshot schedule (%s) resource of storage (%s)-", d.Id(), storageUUID)

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.SuppressHTTPErrorCodes(
		client.DeleteStorageSnapshotSchedule(ctx, storageUUID, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
