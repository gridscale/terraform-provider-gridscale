---
layout: "gridscale"
page_title: "gridscale: storage snapshot schedule"
sidebar_current: "docs-gridscale-datasource-snapshotschedule"
description: |-
  Gets data of a storage snapshot schedule.
---

# gridscale_snapshotschedule

Gets data of a storage snapshot schedule.

## Example Usage

```terraform
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshotschedule" "foo" {
  name = "snapshotschedule"
  storage_uuid = "${gridscale_storage.foo.id}"
  keep_snapshots = 1
  run_interval = 60
  next_runtime = "2025-12-30 15:04:05"
}
data "gridscale_snapshotschedule" "foo" {
	resource_id   = "${gridscale_snapshotschedule.foo.id}"
	storage_uuid   = "${gridscale_storage.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) UUID of the snapshot schedule.

* `storage_uuid` - (Required) UUID of the storage that the snapshot schedule belongs to.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the snapshot schedule.
* `storage_uuid` - UUID of the storage that the snapshot schedule belongs to.
* `status` - The status of the snapshot schedule.
* `name` - The human-readable name of the snapshot schedule.
* `next_runtime` - The date and time that the snapshot schedule will be run.
* `keep_snapshots` - The amount of Snapshots to keep before overwriting the last created Snapshot.
* `run_interval` - The interval at which the schedule will run (in minutes).
* `create_time` - The date and time the snapshot schedule was initially created.
* `change_time` - The date and time of the last snapshot schedule change.
* `labels` - The list of labels.
* `snapshot` - Related snapshots.
    * `name` - Name of the snapshot.
    * `object_uuid` - UUID of the snapshot.
    * `create_time` - The date and time the snapshot was initially created.
