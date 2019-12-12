package gridscale

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/gridscale/gsclient-go"
)

func resourceGridscaleStorageSnapshotSchedule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGridscaleSnapshotScheduleRead,
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
				Description: "Uuid of the storage used to create snapshots",
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
	}
}

func resourceGridscaleSnapshotScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	scheduler, err := client.GetStorageSnapshotSchedule(emptyCtx, storageUUID, d.Id())
	if err != nil {
		if requestError, ok := err.(gsclient.RequestError); ok {
			if requestError.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return err
	}
	props := scheduler.Properties
	d.Set("status", props.Status)
	d.Set("name", props.Name)
	d.Set("next_runtime", props.NextRuntime.String())
	d.Set("keep_snapshots", props.KeepSnapshots)
	d.Set("run_interval", props.RunInterval)
	d.Set("storage_uuid", props.StorageUUID)
	d.Set("create_time", props.CreateTime.String())
	d.Set("change_time", props.ChangeTime.String())

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
		return fmt.Errorf("Error setting snapshots: %v", err)
	}

	//Set labels
	if err = d.Set("labels", props.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %v", err)
	}
	return nil
}
