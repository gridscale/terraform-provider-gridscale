package gridscale

import (
	"fmt"

	"github.com/gridscale/gsclient-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceGridscaleStorageSnapshotSchedule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleSnapshotScheduleRead,

		Schema: map[string]*schema.Schema{
			"resource_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "ID of a resource",
				ValidateFunc: validation.NoZeroValues,
			},
			"storage_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Uuid of the storage used to create snapshots",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the object",
			},
			"next_runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the snapshot schedule will be run",
			},
			"keep_snapshots": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of Snapshots to keep before overwriting the last created Snapshot",
			},
			"run_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The interval at which the schedule will run (in minutes)",
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
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"snapshot": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Related snapshots",
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

func dataSourceGridscaleSnapshotScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)

	storageUUID := d.Get("storage_uuid").(string)
	scheduler, err := client.GetStorageSnapshotSchedule(emptyCtx, storageUUID, id)

	if err == nil {
		props := scheduler.Properties
		d.SetId(props.ObjectUUID)
		if err = d.Set("status", props.Status); err != nil {
			return fmt.Errorf("error setting status: %v", err)
		}
		if err = d.Set("name", props.Name); err != nil {
			return fmt.Errorf("error setting name: %v", err)
		}
		if err = d.Set("next_runtime", props.NextRuntime.Format(timeLayout)); err != nil {
			return fmt.Errorf("error setting next_runtime: %v", err)
		}
		if err = d.Set("keep_snapshots", props.KeepSnapshots); err != nil {
			return fmt.Errorf("error setting keep_snapshots: %v", err)
		}
		if err = d.Set("run_interval", props.RunInterval); err != nil {
			return fmt.Errorf("error setting run_interval: %v", err)
		}
		if err = d.Set("storage_uuid", props.StorageUUID); err != nil {
			return fmt.Errorf("error setting storage_uuid: %v", err)
		}
		if err = d.Set("create_time", props.CreateTime.String()); err != nil {
			return fmt.Errorf("error setting create_time: %v", err)
		}
		if err = d.Set("change_time", props.ChangeTime.String()); err != nil {
			return fmt.Errorf("error setting change_time: %v", err)
		}

		if err = d.Set("labels", props.Labels); err != nil {
			return fmt.Errorf("error setting labels: %v", err)
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
			return fmt.Errorf("error setting snapshots: %v", err)
		}
	}

	return err
}
