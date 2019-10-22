---
layout: "gridscale"
page_title: "gridscale: loadbalancer"
sidebar_current: "docs-gridscale-datasource-loadbalancer"
description: |-
  Gets the id of a loadbalancer.
---

# gridscale_loadbalancer

Get the id of an loabalancer resource.

## Example Usage

```terraform
data "gridscale_loadbalancer" "foo" {
	resource_id   = "xxxx-xxxx-xxxx-xxxx"
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the loadbalancer.
