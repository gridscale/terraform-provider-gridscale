---
layout: "gridscale"
page_title: "gridscale: public network"
sidebar_current: "docs-gridscale-datasource-public-network"
description: |-
  Gets data of a public network.
---

# gridscale_public_network

Get data of the public network. Use this to link your servers to the public network easily.

## Example Usage

Using the public network datasource for the creation of a server:

```terraform
data "gridscale_public_network" "pubnet"{
	project = "default"
}

resource "gridscale_server" "servername"{
	project = data.gridscale_public_network.pubnet.project
	name = "terra-server"
	cores = 2
	memory = 4
	network {
		object_uuid = data.gridscale_public_network.pubnet.id
		bootdevice = true
	}
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.

* `resource_id` - (Required) The UUID of the public network.

## Attributes Reference

The following attributes are exported:

* `project` - The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.
* `name` - The UUID of the network.
* `location_uuid` - The UUID of the location, that helps to identify which datacenter the network belongs to.
* `l2security` - Defines information about MAC spoofing protection.
* `status` - The status of the network.
* `network_type` - The type of the network.
* `location_country` - The human-readable name of the country where the network is located.
* `location_iata` - The IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location where the network is located.
* `delete_block` - Defines if the network is administratively blocked.
* `create_time` - Defines the date and time the network was initially created.
* `change_time` - Defines the date and time of the last network change.
* `labels` - The list of labels.
