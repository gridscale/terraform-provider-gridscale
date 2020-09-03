package gridscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	errHandler "github.com/terraform-providers/terraform-provider-gridscale/gridscale/error-handler"

	"github.com/gridscale/gsclient-go/v3"
)

func resourceGridscaleStorageBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGridscaleBackupScheduleCreate,
		Read:   resourceGridscaleBackupScheduleRead,
		Delete: resourceGridscaleBackupScheduleDelete,
		Update: resourceGridscaleBackupScheduleUpdate,
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
				Required:    true,
				Description: "The date and time that the storage backup schedule will be run. Format: \"2006-01-02 15:04:05\"",
			},
			"next_runtime_computed": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the storage backup schedule will be run. This date and time is computed by gridscale's server.",
			},
			"keep_backups": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The amount of storage backups to keep before overwriting the last created backup",
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
				Description: "Uuid of the storage used to create storage backups",
			},
			"active": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The status of the schedule active or not",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceGridscaleBackupScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("read storage backup schedule (%s) resource of storage (%s)-", d.Id(), storageUUID)
	scheduler, err := client.GetStorageBackupSchedule(context.Background(), storageUUID, d.Id())
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
	if err = d.Set("active", props.Active); err != nil {
		return fmt.Errorf("%s error setting status: %v", errorPrefix, err)
	}
	if err = d.Set("name", props.Name); err != nil {
		return fmt.Errorf("%s error setting name: %v", errorPrefix, err)
	}
	if err = d.Set("next_runtime_computed", props.NextRuntime.Format(timeLayout)); err != nil {
		return fmt.Errorf("%s error setting next_runtime_computed: %v", errorPrefix, err)
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

	//Get storage backups
	backups := make([]interface{}, 0)
	for _, value := range props.Relations.StorageBackups {
		backups = append(backups, map[string]interface{}{
			"name":        value.Name,
			"object_uuid": value.ObjectUUID,
			"create_time": value.CreateTime.String(),
		})
	}
	if err = d.Set("storage_backups", backups); err != nil {
		return fmt.Errorf("%s error setting storage backups: %v", errorPrefix, err)
	}
	return nil
}

func resourceGridscaleBackupScheduleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	requestBody := gsclient.StorageBackupScheduleCreateRequest{
		Name:        d.Get("name").(string),
		RunInterval: d.Get("run_interval").(int),
		KeepBackups: d.Get("keep_backups").(int),
		Active:      d.Get("active").(bool),
	}
	nextRuntime, err := time.Parse(timeLayout, d.Get("next_runtime").(string))
	if err != nil {
		return err
	}
	requestBody.NextRuntime = gsclient.GSTime{Time: nextRuntime}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutCreate))
	defer cancel()
	response, err := client.CreateStorageBackupSchedule(ctx, d.Get("storage_uuid").(string), requestBody)
	if err != nil {
		return err
	}
	d.SetId(response.ObjectUUID)
	log.Printf("The id for storage backup schedule %s has been set to %v", requestBody.Name, response.ObjectUUID)
	return resourceGridscaleBackupScheduleRead(d, meta)
}

func resourceGridscaleBackupScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("update storage backup schedule (%s) resource of storage (%s)-", d.Id(), storageUUID)

	requestBody := gsclient.StorageBackupScheduleUpdateRequest{
		Name:        d.Get("name").(string),
		RunInterval: d.Get("run_interval").(int),
		KeepBackups: d.Get("keep_backups").(int),
	}
	active := d.Get("active").(bool)
	requestBody.Active = &active
	nextRuntime, err := time.Parse(timeLayout, d.Get("next_runtime").(string))
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	requestBody.NextRuntime = &gsclient.GSTime{Time: nextRuntime}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutUpdate))
	defer cancel()
	err = client.UpdateStorageBackupSchedule(ctx, storageUUID, d.Id(), requestBody)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return resourceGridscaleBackupScheduleRead(d, meta)
}

func resourceGridscaleBackupScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)
	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("delete storage backup schedule (%s) resource of storage (%s)-", d.Id(), storageUUID)

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err := errHandler.RemoveErrorContainsHTTPCodes(
		client.DeleteStorageBackupSchedule(ctx, storageUUID, d.Id()),
		http.StatusNotFound,
	)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	return nil
}
