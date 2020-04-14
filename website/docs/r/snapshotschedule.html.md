---
layout: "gridscale"
page_title: "gridscale: storage snapshot schedule"
sidebar_current: "docs-gridscale-resource-snapshotschedule"
description: |-
  Manages a storage snapshot schedule.
---

# gridscale_snapshotschedule

Provides a storage snapshot schedule resource. This can be used to create, modify and delete snapshot schedules.

## Example Usage

```terraform
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshotschedule" "foo" {
  name = "snapshotschedule"
  storage_uuid = gridscale_storage.foo.id
  keep_snapshots = 1
  run_interval = 60
  next_runtime = "2025-12-30 15:04:05"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) UUID of the snapshot schedule.

* `storage_uuid` - (Required) UUID of the storage that the snapshot schedule belongs to.

* `labels` - (Optional) The list of labels.

* `next_runtime` - (Optional) The date and time that the snapshot schedule will be run.

* `keep_snapshots` - (Required) The amount of Snapshots to keep before overwriting the last created Snapshot (>=1).

* `run_interval` - (Required) The interval at which the schedule will run (in minutes, >=60).

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is the value of global `timeout`) Used for Creating resource.
* `update` - (Default value is the value of global `timeout`) Used for Updating resource.
* `delete` - (Default value is the value of global `timeout`) Used for Deleteing resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the snapshot schedule.
* `storage_uuid` - See Argument Reference above.
* `status` - The status of the snapshot schedule.
* `name` - See Argument Reference above.
* `next_runtime` - See Argument Reference above.
* `keep_snapshots` - See Argument Reference above.
* `run_interval` - See Argument Reference above.
* `create_time` - The date and time the snapshot schedule was initially created.
* `change_time` - The date and time of the last snapshot schedule change.
* `labels` - See Argument Reference above.
* `snapshot` - Related snapshots.
    * `name` - Name of the snapshot.
    * `object_uuid` - UUID of the snapshot.
    * `create_time` - The date and time the snapshot was initially created.
