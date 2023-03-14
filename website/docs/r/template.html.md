---
layout: "gridscale"
page_title: "gridscale: template"
sidebar_current: "docs-gridscale-resource-template"
description: |-
  Manages a template in gridscale.
---

# gridscale_template

Provides a template resource. This can be used to create, modify, and delete template.

## Example Usage

The following example shows how one might use this resource to add a template to gridscale:

```terraform
resource "gridscale_storage" "foo" {
  name   = "newname"
  capacity = 1
}

resource "gridscale_snapshot" "foo" {
  name = "newname"
  storage_uuid = gridscale_storage.foo.id
}

resource "gridscale_template" "foo" {
  name   = "newname"
  snapshot_uuid = gridscale_snapshot.foo.id
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The exact name of the template as show in [the page Template](https://my.gridscale.io/Template).

* `snapshot_uuid` - (Required) Snapshot uuid for template.

* `labels` - (Optional) List of labels.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the template.
* `id` - The UUID of the template.
* `location_uuid` - The location this object is placed.
* `location_country` - Two digit country code (ISO 3166-2) of the location where this object is placed.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `status` - Status indicates the status of the object.
* `ostype` - The operating system installed in the template.
* `version` - The version of the template.
* `private` - The object is private, the value will be true. Otherwise the value will be false.
* `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
* `create_time` - The date and time the object was initially created.
* `change_time` - The date and time of the last object change.
* `distro` - The OS distribution that the template contains.
* `description` - Description of the template.
* `usage_in_minutes` - Total minutes the object has been running.
* `capacity` - The capacity of a storage/ISO Image/template/snapshot in GB.
* `current_price` - Defines the price for the current period since the last bill.
* `labels` - List of labels.
