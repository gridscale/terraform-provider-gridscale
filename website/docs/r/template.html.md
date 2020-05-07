---
layout: "gridscale"
page_title: "gridscale: template"
sidebar_current: "docs-gridscale-resource-template"
description: |-
  Manages a template in gridscale.
---

# gridscale_template

Provides a template resource. This can be used to create, modify and delete template.

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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The exact name of the template as show in [the expert panel of gridscale](https://my.gridscale.io/Expert/Template).

* `snapshot_uuid` - (Required) Snapshot uuid for template.

* `labels` - (Optional) List of labels.

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is 0, it means the operation runs without a timeout) Used for Creating resource.
* `update` - (Default value is 0, it means the operation runs without a timeout) Used for Updating resource.
* `delete` - (Default value is 0, it means the operation runs without a timeout) Used for Deleteing resource.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the template.
* `id` - The UUID of the template.
* `location_uuid` - Helps to identify which datacenter an object belongs to.
* `location_country` - Formatted by the 2 digit country code (ISO 3166-2) of the host country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
* `status` - Status indicates the status of the object.
* `ostype` - The operating system installed in the template.
* `version` - The version of the template.
* `private` - The object is private, the value will be true. Otherwise the value will be false.
* `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
* `create_time` - The date and time the object was initially created.
* `change_time` - The date and time of the last object change.
* `distro` - The OS distrobution that the Template contains.
* `description` - Description of the Template.
* `usage_in_minutes` - Total minutes the object has been running.
* `capacity` - The capacity of a storage/ISO Image/template/snapshot in GB.
* `current_price` - Defines the price for the current period since the last bill.
* `labels` - List of labels.
