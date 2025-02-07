---
layout: "gridscale"
page_title: "gridscale: gridscale_k8s"
sidebar_current: "docs-gridscale-resource-k8s"
description: |-
  Get data from a k8s cluster in gridscale.
---

# gridscale_k8s


Get information about a Gridscale Kubernetes cluster.

## Example Usage

```terraform
data "gridscale_k8s" "k8s-example" {
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the Kubernetes cluster.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "45m" - 45 minutes) Used for creating a resource.
* `update` - (Default value is "45m" - 45 minutes) Used for updating a resource.
* `delete` - (Default value is "45m" - 45 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - The human-readable name of the Kubernetes cluster.
* `labels` - See Argument Reference above.
* `kubeconfig` - The kubeconfig file content of the k8s cluster.
* `k8s_private_network_uuid` - Private network UUID which k8s nodes are attached to. It can be used to attach other PaaS/VMs.