---
layout: "gridscale"
page_title: "gridscale: gridscale_network"
sidebar_current: "docs-gridscale-resource-network"
description: |-
  Manages a network in gridscale.
---

# gridscale_network

Provides a network resource. This can be used to create, modify, and delete networks.

## Example Usage

The following example shows how one might use this resource to add a network to gridscale:

```terraform
resource "gridscale_network" "networkname"{
  name = "terraform-network"
  dhcp_active = true
  dhcp_gateway = "192.168.121.1"
  dhcp_dns = "192.168.121.2"
  dhcp_range = "192.168.121.0/27"
  dhcp_reserved_subnet = ["192.168.121.0/31"]
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `l2security` - (Optional) Defines information about MAC spoofing protection (filters layer2 and ARP traffic based on MAC source). It can only be (de-)activated on a private network - the public network always has l2security enabled. It will be true if the network is public, and false if the network is private.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `dhcp_active` - (Optional) Enable DHCP.

* `dhcp_gateway` - (Optional) The general IP Range configured for this network (/24 for private networks). If it is not set, gridscale internal default range is used.

* `dhcp_dns` - (Optional) The IP address reserved and communicated by the dhcp service to be the default gateway.

* `dhcp_range` - (Optional, Computed) DHCP DNS. If it is not set and DHCP is enabled, `dhcp_range` will be set by the backend automatically.

* `dhcp_reserved_subnet` - (Optional) Subrange within the IP range.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `location_uuid` - The location this network is placed. The location of a resource is determined by it's project.
* `l2security` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `dhcp_active` - See Argument Reference above.
* `dhcp_gateway` - See Argument Reference above.
* `dhcp_dns` - See Argument Reference above.
* `dhcp_range` - See Argument Reference above.
* `dhcp_reserved_subnet` - See Argument Reference above.
* `auto_assigned_servers` - A list of server UUIDs with the corresponding IPs that are designated by the DHCP server.
  * `server_uuid` - UUID of the server.
  * `ip` - IP which is assigned to the server.
* `pinned_servers` - A list of server UUIDs with the corresponding IPs that are designated by the user.
  * `server_uuid` - UUID of the server.
  * `ip` - IP which is assigned to the server.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `network_type` - The type of this network, can be mpls, breakout or network.
* `location_country` - Two digit country code (ISO 3166-2) of the location where this object is placed.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `public_net` - Is the network public or not.
* `delete_block` - If deleting this network is allowed.
