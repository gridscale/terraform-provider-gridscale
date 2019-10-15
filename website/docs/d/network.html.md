---
layout: "gridscale"
page_title: "gridscale: network"
sidebar_current: "docs-gridscale-datasource-network"
description: |-
  Gets the id of a network.
---

# gridscale_storage

Get the id of a network resource. This can be used to link networks to a server.

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
		object_uuid = "${data.gridscale_network.networkname.id}"
		bootdevice = true
	}
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the network.
