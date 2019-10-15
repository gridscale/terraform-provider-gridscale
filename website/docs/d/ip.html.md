---
layout: "gridscale"
page_title: "gridscale: ip"
sidebar_current: "docs-gridscale-datasource-ip"
description: |-
  Gets the id of an ip.
---

# gridscale_ip

Get the id of an ip resource. This can be used to link ip addresses to a server.

## Example Usage

Using the ip datasource for the creation of a server:

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
	ipv4 = "${data.gridscale_ipv4.ipv4name.id}"
	ipv6 = "${data.gridscale_ipv6.ipv6name.id}"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the ip.
