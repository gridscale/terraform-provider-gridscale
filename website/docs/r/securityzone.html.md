---
layout: "gridscale"
page_title: "gridscale: securityzone"
sidebar_current: "docs-gridscale-resource-securityzone"
description: |-
  Manages a security zone in gridscale.
---

# gridscale_paas_securityzone

Provides a security zone resource. This can be used to create, modify and delete security zones. 

## Example Usage

The following example shows how one might use this resource to add a security zone to gridscale:

```terraform
resource "gridscale_paas_securityzone" "foo" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
* `location_uuid` - (Optional) Helps to identify which datacenter an object belongs to.

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is the value of global `timeout`) Used for Creating resource.
* `update` - (Default value is the value of global `timeout`) Used for Updating resource.
* `delete` - (Default value is the value of global `timeout`) Used for Deleteing resource.

## Attributes

This resource exports the following attributes:

* `id` - The UUID of the security zone.
* `name` - The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.
* `location_uuid` - Helps to identify which datacenter an object belongs to.
* `location_country` - The human-readable name of the location's country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location.
* `create_time` - Defines the date and time the object was initially created.
* `change_time` - Defines the date and time of the last object change.
* `status` - Status indicates the status of the object.
* `labels` - List of labels.
* `relations` - List of PaaS services' UUIDs relating to the security zone.