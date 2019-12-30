---
layout: "gridscale"
page_title: "gridscale: securityzone"
sidebar_current: "docs-gridscale-datasource-securityzone"
description: |-
  Get a security zone.
---

# gridscale_paas_securityzone

Get a security zone.

## Example Usage

Using the security zone datasource for the creation of a paas:

```terraform
data "gridscale_paas_securityzone" "foo"{
	resource_id = "xxxx-xxxx-xxxx-xxxx"
}


resource "gridscale_paas" "foo"{
	name = "terra-paas-test"
    service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
    security_zone_uuid = data.gridscale_paas_securityzone.foo.id
}
```

## Attributes Reference

The following attributes are exported:

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
