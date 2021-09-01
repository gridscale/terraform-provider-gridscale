---
layout: "gridscale"
page_title: "gridscale: network"
sidebar_current: "docs-gridscale-datasource-network"
description: |-
  Get the data of a network.
---

# gridscale_network

Get data of a network resource. This can be used to link networks to a server.

## Example Usage

Using the network datasource for the creation of a server:

```terraform
data "gridscale_network" "networkname"{
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}

resource "gridscale_server" "servername"{
  name = "terra-server"
  cores = 2
  memory = 4
  network {
    object_uuid = data.gridscale_network.networkname.id
    bootdevice = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the network.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the network.
* `name` - The UUID of the network.
* `location_uuid` - The UUID of the location, that helps to identify which datacenter the network belongs to.
* `l2security` - Defines information about MAC spoofing protection.
* `dhcp_active` - Tell if DHCP is enabled.
* `dhcp_gateway` - The general IP Range configured for this network (/24 for private networks). 
* `dhcp_dns` - The IP address reserved and communicated by the dhcp service to be the default gateway.
* `dhcp_range` -  DHCP DNS.
* `dhcp_reserved_subnet` -  Subrange within the IP range.
* `auto_assigned_servers` - Contains IP addresses of all servers in the network which got a designated IP by the DHCP server.
  * `server_uuid` - UUID of the server.
  * `ip` - IP which is assigned to the server.
* `pinned_servers` - Contains IP addresses of all servers in the network which got a designated IP by the user.
  * `server_uuid` - UUID of the server.
  * `ip` - IP which is assigned to the server.
* `status` - The status of the network.
* `network_type` - The type of the network.
* `location_country` - The human-readable name of the country where the network is located.
* `location_iata` - The IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location where the network is located.
* `delete_block` - Defines if the network is administratively blocked.
* `create_time` - Defines the date and time the network was initially created.
* `change_time` - Defines the date and time of the last network change.
* `labels` - The list of labels.
