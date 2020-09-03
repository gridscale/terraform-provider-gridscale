---
layout: "gridscale"
page_title: "gridscale: storage backup schedule"
sidebar_current: "docs-gridscale-resource-backupschedule"
description: |-
  Manages a storage backup schedule.
---

# gridscale_backupschedule

Provides a storage backup schedule resource. This can be used to create, modify and delete backup schedules.

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
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) UUID of the backup schedule.

* `storage_uuid` - (Required) UUID of the storage that the backup schedule belongs to.

* `active` - (Required) The status of the schedule active or not.

* `next_runtime` - (Required) The date and time that the backup schedule will be run.

* `keep_backups` - (Required) The amount of Snapshots to keep before overwriting the last created Snapshot (>=1).

* `run_interval` - (Required) The interval at which the schedule will run (in minutes, >=60).

## Timeouts

Timeouts configuration options (in seconds):
More info: https://www.terraform.io/docs/configuration/resources.html#operation-timeouts

* `create` - (Default value is "5m" - 5 minutes) Used for Creating resource.
* `update` - (Default value is "5m" - 5 minutes) Used for Updating resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for Deleteing resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the backup schedule.
* `storage_uuid` - See Argument Reference above.
* `status` - The status of the backup schedule.
* `active` - See Argument Reference above.
* `name` - See Argument Reference above.
* `next_runtime` - See Argument Reference above.
* `next_runtime_computed` - The date and time that the backup schedule will be run. This date and time is computed by gridscale's server.
* `keep_backups` - See Argument Reference above.
* `run_interval` - See Argument Reference above.
* `create_time` - The date and time the backup schedule was initially created.
* `change_time` - The date and time of the last backup schedule change.
* `backup` - Related backups.
    * `name` - Name of the backup.
    * `object_uuid` - UUID of the backup.
    * `create_time` - The date and time the backup was initially created.
