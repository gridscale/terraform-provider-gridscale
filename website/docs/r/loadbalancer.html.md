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
	listen_ipv4_uuid = gridscale_ipv4.lb.id
	listen_ipv6_uuid = gridscale_ipv6.lb.id
	labels = []
	backend_server {
		weight = 100
		host   = gridscale_ipv4.server.ip
	}
	forwarding_rule {
		listen_port =  80
		mode        =  "http"
		target_port =  80
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `redirect_http_to_https` - (Required) Whether the loadbalancer is forced to redirect requests from HTTP to HTTPS.

* `listen_ipv4_uuid` - (Required) The UUID of the IPv4 address the loadbalancer will listen to for incoming requests.

* `listen_ipv6_uuid` - (Required) The UUID of the IPv6 address the loadbalancer will listen to for incoming requests.

* `algorithm` - (Required) The algorithm used to process requests. Accepted values: roundrobin/leastconn.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is the value of global `timeout`) Used for Creating resource.
* `update` - (Default value is the value of global `timeout`) Used for Updating resource.
* `delete` - (Default value is the value of global `timeout`) Used for Deleteing resource.

## Attributes

This resource exports the following attributes:

* `id` - The UUID of the loadbalancer.
* `location_uuid` - Helps to identify which datacenter an object belongs to. The location of the resource depends on the location of the project.
* `name` - The human-readable name of the loadbalancer.
* `algorithm` - The algorithm used to process requests.
* `status` - The status of the loadbalancer.
* `redirect_http_to_https` - Whether the Load balancer is forced to redirect requests from HTTP to HTTPS.
* `listen_ipv4_uuid` - The UUID of the IPv4 address the loadbalancer will listen to for incoming requests.
* `listen_ipv6_uuid` - The UUID of the IPv6 address the loadbalancer will listen to for incoming requests.
* `forwarding_rule` - The forwarding rules of the loadbalancer.
* `backend_server` - The servers that the loadbalancer can communicate with.
* `labels` - The list of labels.
