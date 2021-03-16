---
layout: "gridscale"
page_title: "gridscale: gridscale_k8s"
sidebar_current: "docs-gridscale-resource-k8s"
description: |-
  Manages a k8s cluster in gridscale.
---

# gridscale_k8s


Provides a k8s cluster resource. This can be used to create, modify, and delete k8s cluster resource.

## Example Usage

The following example shows how one might use this resource to add a k8s cluster to gridscale:

```terraform
resource "gridscale_k8s" "k8s-test" {
  name   = "test"
  k8s_release = "1.19"
  node_pool {
	name = "my_node_pool"
    node_count = 2
    cores = 1
    memory = 2
    storage = 10
    storage_type = "storage_insane"
  }
 }

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `security_zone_uuid` - (Optional, Computed, ForceNew) Security zone UUID linked to k8s service. If `security_zone_uuid` is not set, the default security zone will be created (if it doesn't exist) and linked to the k8s service.

* `k8s_release` - (Required) Release number of k8s service. Define which release of k8s service will be used to create the k8s cluster.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `node_pool` - (Required) Node pool's specification. **NOTE**: The node pool's specification is not yet mutable.
    * `name` - Name of the node pool.
    * `node_count` - Number of worker nodes.
    * `cores` - Cores per worker node.
    * `memory` - Memory per worker node (in GiB).
    * `storage` - Storage per worker node (in GiB).
    * `storage_type` - Storage type (one of storage, storage_high, storage_insane).

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "15m" - 15 minutes) Used for creating a resource.
* `update` - (Default value is "15m" - 15 minutes) Used for updating a resource.
* `delete` - (Default value is "15m" - 15 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `security_zone_uuid` - See Argument Reference above.
* `k8s_release` - See Argument Reference above.
* `k8s_release_computed` - Release number of k8s service. The `k8s_release_computed` will be different from `k8s_release`, when `k8s_release` is updated outside of terraform.
* `labels` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone, which is linked to the k8s cluster.
* `node_pool` - See Argument Reference above.
    * `name` - See Argument Reference above.
    * `node_count` - See Argument Reference above.
    * `cores` - See Argument Reference above.
    * `memory` - See Argument Reference above.
    * `storage` - See Argument Reference above.
    * `storage_type` - See Argument Reference above.
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
* `current_price` - The price for the current period since the last bill.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `status` - status indicates the status of the object.
