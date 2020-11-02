---
layout: "gridscale"
page_title: "gridscale: storage backup schedule"
sidebar_current: "docs-gridscale-datasource-backupschedule"
description: |-
  Gets data of a storage backup schedule.
---

# gridscale_backupschedule

Gets data of a storage backup schedule.

## Example Usage

```terraform
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_backupschedule" "foo" {
  name = "backupschedule"
  storage_uuid = gridscale_storage.foo.id
  keep_backups = 1
  run_interval = 60
  next_runtime = "2025-12-30 15:04:05"
  active       = true
}
data "gridscale_backupschedule" "foo" {
  resource_id   = gridscale_backupschedule.foo.id
  storage_uuid   = gridscale_storage.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) UUID of the backup schedule.

* `storage_uuid` - (Required) UUID of the storage that the backup schedule belongs to.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the backup schedule.
* `storage_uuid` - UUID of the storage that the backup schedule belongs to.
* `status` - The status of the backup schedule.
* `active` - The status of the schedule active or not.
* `name` - The human-readable name of the backup schedule.
* `next_runtime` - The date and time that the backup schedule will be run.
* `keep_backups` - The amount of Snapshots to keep before overwriting the last created Snapshot.
* `run_interval` - The interval at which the schedule will run (in minutes).
* `create_time` - The date and time the backup schedule was initially created.
* `change_time` - The date and time of the last backup schedule change.
* `storage_backups` - Related backups.
  * `name` - Name of the backup.
  * `object_uuid` - UUID of the backup.
  * `create_time` - The date and time the backup was initially created.
