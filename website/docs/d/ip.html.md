---
layout: "gridscale"
page_title: "gridscale: ip"
sidebar_current: "docs-gridscale-datasource-ip"
description: |-
  Gets data of an IP address.
---

# gridscale_ip

Get data of an IP address resource. This can be used to assign IP addresses to a server.

## Example Usage

Using ip datasource for the creation of a server:

```terraform
data "gridscale_ipv4" "ipv4name"{
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}

data "gridscale_ipv6" "ipv6name"{
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}

resource "gridscale_server" "servername"{
  name = "terra-server"
  cores = 2
  memory = 4
  ipv4 = data.gridscale_ipv4.ipv4name.id
  ipv6 = data.gridscale_ipv6.ipv6name.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the IP.
* `ip` - Defines the IP address (v4 or v6).
* `prefix` - The prefix of the IP or subnet.
* `location_uuid` - The UUID of the location, that helps to identify which datacenter an object belongs to.
* `failover` - failover mode of this ip. If true, then this IP is no longer available for DHCP and can no longer be related to any server.
* `status` - The status of the IP.
* `reverse_dns` - The reverse DNS of the IP.
* `location_iata` - The IATA airport code, which works as a location identifier.
* `location_country` - The human-readable name of the country of the IP.
* `location_name` - The human-readable name of the location of the IP.
* `create_time` - The date and time the IP was initially created.
* `change_time` - The date and time of the last ip change.
* `delete_block` - Defines if the IP is administratively blocked.
* `usage_in_minutes` - Total minutes the IP has been running.
* `current_price` - The price for the current period since the last bill.
* `labels` - The list of labels.
