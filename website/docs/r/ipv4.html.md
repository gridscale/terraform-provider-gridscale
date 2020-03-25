---
layout: "gridscale"
page_title: "gridscale: gridscale_ipv4"
sidebar_current: "docs-gridscale-resource-ipv4"
description: |-
  Manages an IPv4 address in gridscale.
---

# gridscale_ipv4

Provides an IPv4 address resource. This can be used to create, modify and delete IPv4 addresses.

## Example Usage

The following example shows how one might use this resource to add an IPv4 address to gridscale:

```terraform
resource "gridscale_ipv4" "terra-ipv4-test" {
  project = "default"
	name = "terra-test"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.

* `name` - (Optional) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `failover` - (Optional) Sets failover mode for this IP. If true, then this IP is no longer available for DHCP and can no longer be related to any server.

* `reverse_dns` - (Optional) Defines the reverse DNS entry for the IP Address (PTR Resource Record).

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

## Attributes

This resource exports the following attributes:

* `project` - The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.
* `id` - The UUID of the IPv4.
* `name` - See Argument Reference above.
* `location_uuid` - See Argument Reference above.
* `failover` - See Argument Reference above.
* `reverse_dns` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `ip` - Defines the IP Address.
* `prefix` - The network address and the subnet.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `location_country` - Formatted by the 2 digit country code (ISO 3166-2) of the host country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `delete_block` - Defines if the object is administratively blocked. If true, it can not be deleted by the user.
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
* `current_price` - The price for the current period since the last bill.
