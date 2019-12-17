---
layout: "gridscale"
page_title: "gridscale: public network"
sidebar_current: "docs-gridscale-datasource-public-network"
description: |-
  Gets the id of a public network.
---

# gridscale_public_network

Get date of the public network. This can be used to link public network to a server.

## Example Usage

Using the public network datasource for the creation of a server:

```terraform
data "gridscale_public_network" "pubnet"{
}

resource "gridscale_server" "servername"{
	name = "terra-server"
	cores = 2
	memory = 4
	network {
		object_uuid = "${data.gridscale_public_network.pubnet.id}"
		bootdevice = true
	}
}
```

## Attributes Reference

The following attributes are exported:

* `name` - The UUID of the network.
* `location_uuid` - The UUID of the location, that helps to identify which datacenter the network belongs to.
* `l2security` - Defines information about MAC spoofing protection.
* `status` - The status of the network.
* `network_type` - The type of the network.
* `location_country` - The human-readable name of the country where the network locates.
* `location_iata` - The IATA airport code, which works as a location identifier.
* `location_name` - The uman-readable name of the location where the network locates.
* `delete_block` - Defines if the network is administratively blocked.
* `create_time` - Defines the date and time the network was initially created.
* `change_time` - Defines the date and time of the last network change.
* `labels` - The list of labels.
