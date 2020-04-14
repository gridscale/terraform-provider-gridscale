---
layout: "gridscale"
page_title: "gridscale: storage snapshot"
sidebar_current: "docs-gridscale-resource-snapshot"
description: |-
  Manages a storage snapshot in gridscale.
---

# gridscale_snapshot

Provides a storage snapshot resource. This can be used to create, modify and delete storage snapshots.

## Example Usage

```terraform
resource "gridscale_storage" "foo" {
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  name = "snapshot"
  storage_uuid = gridscale_storage.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the snapshot.

* `storage_uuid` - (Required) UUID of the storage used to create this snapshot.

* `labels` - (Optional) The list of labels.

* `rollback` - (Optional) Returns a storage to the state of the selected Snapshot. 

    * `id` - (Required) ID of the rollback request. It can be any string value. Each rollback request has to have a UNIQUE id. 

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is the value of global `timeout`) Used for Creating resource.
* `update` - (Default value is the value of global `timeout`) Used for Updating resource.
* `delete` - (Default value is the value of global `timeout`) Used for Deleteing resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the snapshot.
* `storage_uuid` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - The status of the snapshot.
* `location_uuid` - The UUID of the location, that helps to identify which datacenter an object belongs to.
* `location_iata` - The IATA airport code, which works as a location identifier.
* `location_country` - The human-readable name of the country of the snapshot.
* `location_name` - The human-readable name of the location of the snapshot.
* `create_time` - The date and time the ip was initially created.
* `change_time` - The date and time of the last snapshot change.
* `usage_in_minutes` - Total minutes the ip has been running.
* `current_price` - The price for the current period since the last bill.
* `capacity` - The capacity of the snapshot in GB.
* `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
* `rollback` - See Argument Reference above.
    * `id` - See Argument Reference above. 
    * `rollback_time` - The time when rollback request is fulfilled.
    * `status` - Status of the rollback request.
* `labels` - See Argument Reference above.
