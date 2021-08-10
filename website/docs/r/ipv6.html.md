---
layout: "gridscale"
page_title: "gridscale: gridscale_ipv6"
sidebar_current: "docs-gridscale-resource-ipv6"
description: |-
  Manages an IPv6 address in gridscale.
---

# gridscale_ipv6

Provides an IPv6 address resource. This can be used to create, modify, and delete IPv6 addresses.

## Example Usage

The following example shows how one might use this resource to add an IPv6 address to gridscale:

```terraform
resource "gridscale_ipv6" "terra-ipv6-test" {
  name = "terra-test"
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `failover` - (Optional) Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.

* `reverse_dns` - (Optional) Defines the reverse DNS entry for the IP address (PTR Resource Record).

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
* `location_uuid` - The location this resource is placed. The location of a resource is determined by it's project.
* `failover` - See Argument Reference above.
* `reverse_dns` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `ip` - Defines the IP address.
* `prefix` - The network address and the subnet.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `location_country` - Two digit country code (ISO 3166-2) of the location where this object is placed.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `delete_block` - Defines if the object is administratively blocked. If true, it can not be deleted by the user.
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
