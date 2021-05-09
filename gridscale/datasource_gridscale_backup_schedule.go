package gridscale

import (
	"context"
	"fmt"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceGridscaleStorageBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleBackupScheduleRead,

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
				Description: "UUID of the storage used to create backups",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The human-readable name of the object",
			},
			"next_runtime": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup schedule will be run",
			},
			"keep_backups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of Backups to keep before overwriting the last created Backup",
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
			"active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The status of the schedule active or not",
			},
			"storage_backups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Related backups",
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

func dataSourceGridscaleBackupScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	id := d.Get("resource_id").(string)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("read backup schedule (%s) datasource of storage (%s) -", id, storageUUID)

	scheduler, err := client.GetStorageBackupSchedule(context.Background(), storageUUID, id)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}

	props := scheduler.Properties
	d.SetId(props.ObjectUUID)
	if err = d.Set("status", props.Status); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("active", props.Active); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("next_runtime", props.NextRuntime.Format(timeLayout)); err != nil {
		return fmt.Errorf("%s error setting next_runtime: %v", errorPrefix, err)
	}
	if err = d.Set("keep_backups", props.KeepBackups); err != nil {
		return fmt.Errorf("%s error setting keep_backups: %v", errorPrefix, err)
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

	//Get backups
	backups := make([]interface{}, 0)
	for _, value := range props.Relations.StorageBackups {
		backups = append(backups, map[string]interface{}{
			"name":        value.Name,
			"object_uuid": value.ObjectUUID,
			"create_time": value.CreateTime.String(),
		})
	}
	if err = d.Set("storage_backups", backups); err != nil {
		return fmt.Errorf("%s error setting backups: %v", errorPrefix, err)
	}

	return nil
}
