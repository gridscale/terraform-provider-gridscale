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
  release = "1.21" # instead, gsk_version can be set.
  node_pool {
  name = "my_node_pool"
    node_count = 2
    cores = 1
    memory = 2
    storage = 10
    storage_type = "storage_insane"
    rocket_storage = 500
  }
 }

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `security_zone_uuid` -  *DEPRECATED* (Optional, Forcenew) Security zone UUID linked to the Kubernetes resource. If `security_zone_uuid` is not set, the default security zone will be created (if it doesn't exist) and linked. A change of this argument necessitates the re-creation of the resource.

* `gsk_version` - (Optional) The gridscale's Kubernetes version of this instance (e.g. "1.21.14-gs1"). Define which gridscale k8s version will be used to create the cluster. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available gridscale k8s version. **NOTE**: Either `gsk_version` or `release` is set at a time.

* `release` - (Optional) The Kubernetes release of this instance. Define which release will be used to create the cluster. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available releases. **NOTE**: Either `gsk_version` or `release` is set at a time.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `node_pool` - (Required) Node pool's specification. **NOTE**: The node pool's specification is not yet mutable (except `node_count`).
    * `name` - (Immutable) Name of the node pool.
    * `node_count` - Number of worker nodes.
    * `cores` - (Immutable) Cores per worker node.
    * `memory` - (Immutable) Memory per worker node (in GiB).
    * `storage` - (Immutable) Storage per worker node (in GiB).
    * `storage_type` - (Immutable) Storage type (one of storage, storage_high, storage_insane).
    * `rocket_storage` - Rocket storage per worker node (in GiB).
    * `surge_node` - Enable surge node to avoid resources shortage during the cluster upgrade (Default: true).
    * `cluster_cidr` - The cluster CIDR that will be used to generate the CIDR of nodes, services, and pods. The allowed CIDR prefix length is /16. If the cluster CIDR is not set, the cluster will use "10.244.0.0/16" as it default (even though the `cluster_cidr` in the k8s resource is empty).

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "45m" - 45 minutes) Used for creating a resource.
* `update` - (Default value is "45m" - 45 minutes) Used for updating a resource.
* `delete` - (Default value is "45m" - 45 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `security_zone_uuid` - See Argument Reference above.
* `release` - See Argument Reference above.
* `gsk_version` - See Argument Reference above.
* `service_template_uuid` - PaaS service template that k8s service uses. The `service_template_uuid` may not relate to `release`, if `service_template_uuid`/`release` is updated outside of terraform (e.g. the k8s service is upgraded by gridscale staffs).
* `service_template_category` - The template service's category used to create the service.
* `labels` - See Argument Reference above.
* `kubeconfig` - The kubeconfig file content of the k8s cluster.
* `network_uuid` - *DEPRECATED*  Network UUID containing security zone, which is linked to the k8s cluster.
* `k8s_private_network_uuid` - Private network UUID which k8s nodes are attached to. It can be used to attach other PaaS/VMs.
* `node_pool` - See Argument Reference above.
    * `name` - See Argument Reference above.
    * `node_count` - See Argument Reference above.
    * `cores` - See Argument Reference above.
    * `memory` - See Argument Reference above.
    * `storage` - See Argument Reference above.
    * `storage_type` - See Argument Reference above.
    * `rocket_storage` - See Argument Reference above.
    * `surge_node` - See Argument Reference above.
    * `cluster_cidr` - See Argument Reference above.
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `status` - status indicates the status of the object.
