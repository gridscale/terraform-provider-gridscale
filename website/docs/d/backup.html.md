---
layout: "gridscale"
page_title: "gridscale: storage backup list"
sidebar_current: "docs-gridscale-datasource-backup-list"
description: |-
  Gets a backup list of a specific storage.
---

# gridscale_backup_list

Gets a backup list of a specific storage.

## Example Usage

```terraform
data "gridscale_backup_list" "foo" {
  	storage_uuid = "XXXX-XXXX-XXXX-XXXX"
}
```

## Argument Reference

The following arguments are supported:

* `storage_uuid` - (Required) UUID of the storage that the backups belong to.

## Attributes Reference

The following attributes are exported:

* `storage_uuid` - UUID of the storage that the backups belong to.
* `storage_backups` - Backups of the given storage.
    * `name` - Name of the backup.
    * `object_uuid` - UUID of the backup.
    * `create_time` - The date and time the backup was initially created.
    * `capacity` - The size of a backup in GB.
