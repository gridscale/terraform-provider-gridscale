---
layout: "gridscale"
page_title: "gridscale: securityzone"
sidebar_current: "docs-gridscale-resource-securityzone"
description: |-
  Manages a security zone in gridscale.
---

# gridscale_paas_securityzone

Provides a security zone resource. This can be used to create, modify, and delete security zones.

## Example Usage

The following example shows how one might use this resource to add a security zone to gridscale:

```terraform
resource "gridscale_paas_securityzone" "foo" {
  name = "test"
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `location_uuid` - (Optional) The location this object is placed.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `id` - The UUID of the security zone.
* `name` - The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `location_uuid` - The location this object is placed.
* `location_country` - The human-readable name of the location's country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location.
* `create_time` - Defines the date and time the object was initially created.
* `change_time` - Defines the date and time of the last object change.
* `status` - Status indicates the status of the object.
* `labels` - List of labels.
* `relations` - List of PaaS services' UUIDs relating to the security zone.
