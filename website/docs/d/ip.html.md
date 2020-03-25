---
layout: "gridscale"
page_title: "gridscale: ip"
sidebar_current: "docs-gridscale-datasource-ip"
description: |-
  Gets data of an IP address.
---

# gridscale_ip

Get data of an IP address resource. This can be used to link ip addresses to a server.

## Example Usage

Using ip datasource for the creation of a server:

```terraform
data "gridscale_ipv4" "ipv4name"{
  	project = "default"
	resource_id = "xxxx-xxxx-xxxx-xxxx"
}

data "gridscale_ipv6" "ipv6name"{
	project = "default"
	resource_id = "xxxx-xxxx-xxxx-xxxx"
}

resource "gridscale_server" "servername"{
	project = "default"
	name = "terra-server"
	cores = 2
	memory = 4
	ipv4 = data.gridscale_ipv4.ipv4name.id
	ipv6 = data.gridscale_ipv6.ipv6name.id
}
```
## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.

* `resource_id` - (Required) The UUID of the IP address.

## Attributes Reference

The following attributes are exported:

* `project` - The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.
* `id` - The UUID of the ip.
* `ip` - Defines the IP Address (v4 or v6) the ip.
* `prefix` - The IP prefix of the ip.
* `location_uuid` - The UUID of the location, that helps to identify which datacenter an object belongs to.
* `failover` - failover mode of this ip. If true, then this IP is no longer available for DHCP and can no longer be related to any server..
* `status` - The status of the ip.
* `reverse_dns` - The reverse DNS of the ip.
* `location_iata` - The IATA airport code, which works as a location identifier.
* `location_country` - The human-readable name of the country of the ip.
* `location_name` - The human-readable name of the location of the ip.
* `create_time` - The date and time the ip was initially created.
* `change_time` - The date and time of the last ip change.
* `delete_block` - Defines if the ip is administratively blocked.
* `usage_in_minutes` - Total minutes the ip has been running.
* `current_price` - The price for the current period since the last bill.
* `labels` - The list of labels.
