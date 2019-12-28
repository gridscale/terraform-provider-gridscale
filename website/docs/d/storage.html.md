---
layout: "gridscale"
page_title: "gridscale: storage"
sidebar_current: "docs-gridscale-datasource-storage"
description: |-
  Gets the id of a storage.
---

# gridscale_storage

Get the id of a storage resource. This can be used to link storages to a server.

## Example Usage

Using the storage datasource for the creation of a server:

```terraform
data "gridscale_storage" "storagename"{
	resource_id = "xxxx-xxxx-xxxx-xxxx"
}

resource "gridscale_server" "servername"{
	name = "terra-server"
	cores = 2
	memory = 4
	storage {
		object_uuid = data.gridscale_storage.storagename.id
		bootdevice = true
	}
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the storage.
* `change_time` - Defines the date and time of the last storage change.
* `location_iata` - The IATA airport code of the location where storage locates.
* `status` - The status of the storage.
* `license_product_no` - The license key (e.g. Windows Servers), if the template used by the storage requires.
* `location_country` - The human-readable name of the country where the storage locates.
* `usage_in_minutes` - Total minutes the the storage has been running.
* `last_used_template` - The UUID of the last used template on the storage.
* `current_price` - The price for the current period since the last bill.
* `capacity` - The capacity (GB) of the storage.
* `location_uuid` - The UUID of the location where the storage locates.
* `storage_type` - The type of the storage.
* `parent_uuid` - The UUID of the parent of the storage.
* `name` - The human-readable name of the storage.
* `location_name` - The human-readable name of the location where the storage locates.
* `create_time` - Defines the date and time the storage was initially created.
* `labels` - The list of labels.
