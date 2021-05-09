---
layout: "gridscale"
page_title: "gridscale: gridscale_storage_clone"
sidebar_current: "docs-gridscale-resource-storage-clone"
description: |-
  Make a clone of an existing storage instance.
---

# gridscale_storage_clone

Clone a storage instance. This can be used to create, modify, and delete the storage clones.

## Example Usage

The following example shows how to clone a storage instance in gridscale:

```terraform
resource "gridscale_storage" "storage-john"{
  name = "john's storage"
  capacity = 10
  storage_type = "storage_high"
  timeouts {
    create="10m"
  }
}

resource "gridscale_storage_clone" "storage-clone-john"{
  source_storage_id = gridscale_storage.storage-clone-john.id
  name = "john's storage clone"
  timeouts {
    create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `source_storage_id` - (Required) The ID of a storage instance which will be cloned.

* `name` - (Optional) The default value is inherited from the source storage instance. A desired name is possible. The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `capacity` - (Optional) The default value is inherited from the source storage instance. A desired capacity is possible. Required (integer - minimum: 1 - maximum: 4096).

* `storage_type` - (Optional) The default value is inherited from the source storage instance. A desired storage type is possible. (one of storage, storage_high, storage_insane).

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `capacity` - See Argument Reference above.
* `storage_type` - See Argument Reference above.
* `location_uuid` - The location this resource is placed. The location of a resource is determined by it's project.
* `labels` - See Argument Reference above.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `location_country` - Two digit country code (ISO 3166-2) of the location where this object is placed.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
* `last_used_template` - Indicates the UUID of the last used template on this storage (inherited from snapshots).
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
* `current_price` - The price for the current period since the last bill.
