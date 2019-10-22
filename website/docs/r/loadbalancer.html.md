---
layout: "gridscale"
page_title: "gridscale: loadbalancer"
sidebar_current: "docs-gridscale-resource-loadbalancer"
description: |-
  Manage a loadbalancer in gridscale.
---

# gridscale_loadbalancer

Provides a loadbalancer resource. This can be used to create, modify and delete loadbalancers.

## Example Usage

```terraform
resource "gridscale_loadbalancer" "foo" {
	name   = "%s"
	algorithm = "%s"
	redirect_http_to_https = false
	listen_ipv4_uuid = "${gridscale_ipv4.lb.id}"
	listen_ipv6_uuid = "${gridscale_ipv6.lb.id}"
	labels = []
	backend_server {
		weight = 100
		host   = "${gridscale_ipv4.server.ip}"
	}
	forwarding_rule {
		listen_port =  80
		mode        =  "http"
		target_port =  80
	}
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the loadbalancer.
