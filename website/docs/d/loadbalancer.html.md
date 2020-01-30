---
layout: "gridscale"
page_title: "gridscale: loadbalancer"
sidebar_current: "docs-gridscale-datasource-loadbalancer"
description: |-
  Get the data of a Load Balancer.
---

# gridscale_loadbalancer

Get the data of a Load Balancer.

## Example Usage

```terraform
data "gridscale_loadbalancer" "foo" {
	resource_id   = "xxxx-xxxx-xxxx-xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the loadbalancer.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the loadbalancer.
* `name` - The human-readable name of the loadbalancer.
* `algorithm` - The algorithm used to process requests.
* `status` - The status of the loadbalancer.
* `redirect_http_to_https` - Whether the Load balancer is forced to redirect requests from HTTP to HTTPS.
* `listen_ipv4_uuid` - The UUID of the IPv4 address the loadbalancer will listen to for incoming requests.
* `listen_ipv6_uuid` - The UUID of the IPv6 address the loadbalancer will listen to for incoming requests.
* `forwarding_rule` - The forwarding rules of the loadbalancer.
* `backend_server` - The servers that the loadbalancer can communicate with.
* `labels` - The list of labels.
