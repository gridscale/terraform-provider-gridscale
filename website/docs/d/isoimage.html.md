---
layout: "gridscale"
page_title: "gridscale: ISO-Image"
sidebar_current: "docs-gridscale-datasource-isoimage"
description: |-
  Gets data of an ISO-Image by its UUID.
---

# gridscale_isoimage

Get data of an ISO-Image by its UUID.

## Example Usage

```terraform
resource "gridscale_isoimage" "foo" {
  name   = "name"
  source_url = "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso"
}

data "gridscale_isoimage" "foo" {
	resource_id   = gridscale_isoimage.foo.id
}
```


## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the ISO-Image.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the ISO-Image.
* `source_url` - Contains the source URL of the ISO-Image that it was originally fetched from.
* `server` - The information about servers which are related to this ISO-Image.
    * `object_uuid` - The object UUID or id of the server.
    * `object_name` - Name of the server.
    * `create_time` - The date and time the object was initially created.
    * `bootdevice` - True if the ISO-Image is a boot device of this server.
* `id` - The UUID of the ISO-Image.
* `location_uuid` - Helps to identify which datacenter an object belongs to.
* `location_country` - Formatted by the 2 digit country code (ISO 3166-2) of the host country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
* `status` - Status indicates the status of the object.
* `version` - The version of the ISO-Image.
* `private` - The object is private, the value will be true. Otherwise the value will be false.
* `create_time` - The date and time the object was initially created.
* `change_time` - The date and time of the last object change.
* `description` - Description of the Template.
* `usage_in_minutes` - Total minutes the object has been running.
* `capacity` - The capacity of a storage/ISO-Image/ISO-Image/snapshot in GB.
* `current_price` - Defines the price for the current period since the last bill.
* `labels` - List of labels.
