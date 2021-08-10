---
layout: "gridscale"
page_title: "gridscale: ISO Image"
sidebar_current: "docs-gridscale-resource-isoimage"
description: |-
  Manages an ISO Image in Gridscale.
---

# gridscale_isoimage

Provides an ISO Image resource. This can be used to create, modify, and delete ISO Images.

## Example Usage

```terraform
resource "gridscale_isoimage" "foo" {
  name   = "newname"
  source_url = "http://tinycorelinux.net/10.x/x86/release/TinyCore-current.iso"
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `source_url` - (Required) Contains the source URL of the ISO Image that it was originally fetched from.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the ISO Image.
* `source_url` - Contains the source URL of the ISO Image that it was originally fetched from.
* `server` - The information about servers which are related to this ISO Image.
  * `object_uuid` - The object UUID or id of the server.
  * `object_name` - Name of the server.
  * `create_time` - The date and time the object was initially created.
  * `bootdevice` - True if the ISO Image is a boot device of this server.
* `id` - The UUID of the ISO Image.
* `location_uuid` - The location this object is placed.
* `location_country` - Two digit country code (ISO 3166-2) of the location where this object is placed.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `status` - Status indicates the status of the object.
* `version` - The version of the ISO Image.
* `private` - The object is private, the value will be true. Otherwise the value will be false.
* `create_time` - The date and time the object was initially created.
* `change_time` - The date and time of the last object change.
* `description` - Description of the template.
* `usage_in_minutes` - Total minutes the object has been running.
* `capacity` - The capacity of a storage/ISO Image/ISO Image/snapshot in GB.
* `labels` - List of labels.
