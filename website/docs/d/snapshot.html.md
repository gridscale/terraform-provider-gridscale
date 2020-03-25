---
layout: "gridscale"
page_title: "gridscale: storage snapshot"
sidebar_current: "docs-gridscale-datasource-snapshot"
description: |-
  Gets data of a storage snapshot.
---

# gridscale_snapshot

Get data of a storage snapshot resource.

## Example Usage

```terraform
resource "gridscale_storage" "foo" {
  project = "default"
  name   = "storage"
  capacity = 1
}
resource "gridscale_snapshot" "foo" {
  project = gridscale_storage.foo.project
  name = "snapshot"
  storage_uuid = gridscale_storage.foo.id
}

data "gridscale_snapshot" "foo" {
	  project   = gridscale_snapshot.foo.project
	  resource_id   = gridscale_snapshot.foo.id
  	storage_uuid = gridscale_storage.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.

* `resource_id` - (Required) ID of a resource (UUID of snapshot).

* `storage_uuid` - (Required) UUID of the storage used to create this snapshot.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the snapshot.
* `storage_uuid` - The UUID of the storage used to create this snapshot.
* `name` - The name of the snapshot.
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
* `labels` - The list of labels.
