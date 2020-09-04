package gridscale

import (
	"context"
	"fmt"
	"sort"

	"github.com/gridscale/gsclient-go/v3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGridscaleStorageBackupList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGridscaleBackupListRead,
		Schema: map[string]*schema.Schema{
			"storage_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the storage",
			},
			"storage_backups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Storage's backups. The order is based on their created time. E.g: Latest backup is always the first backup in the list",
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
						"capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGridscaleBackupListRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gsclient.Client)

	storageUUID := d.Get("storage_uuid").(string)
	errorPrefix := fmt.Sprintf("read backups datasource of storage (%s)-", storageUUID)

	backupList, err := client.GetStorageBackupList(context.Background(), storageUUID)
	if err != nil {
		return fmt.Errorf("%s error: %v", errorPrefix, err)
	}
	d.SetId(storageUUID)
	//Sort the backup List by create_time
	sort.SliceStable(backupList, func(i, j int) bool {
		iCreateTime := backupList[i].Properties.CreateTime
		jCreateTime := backupList[j].Properties.CreateTime.Time
		return iCreateTime.After(jCreateTime)
	})
	//Get storage backups
	backups := make([]interface{}, 0)
	for _, value := range backupList {
		prop := value.Properties
		backups = append(backups, map[string]interface{}{
			"name":        prop.Name,
			"object_uuid": prop.ObjectUUID,
			"create_time": prop.CreateTime.String(),
			"capacity":    prop.Capacity,
		})
	}

	if err = d.Set("storage_backups", backups); err != nil {
		return fmt.Errorf("%s error setting storage backups: %v", errorPrefix, err)
	}
	return nil
}
