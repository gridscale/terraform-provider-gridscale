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
	release = "1.30" # instead, gsk_version can be set.

	node_pool {
    name = "pool-0"
    node_count = 2
    cores = 2
    memory = 4
    storage = 30
    storage_type = "storage_insane"
  }

	node_pool {
    name = "pool-1"
    node_count = 3
    cores = 1
    memory = 3
    storage = 30
    storage_type = "storage_insane"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `security_zone_uuid` -  *DEPRECATED* (Optional, Forcenew) Security zone UUID linked to the Kubernetes resource. If `security_zone_uuid` is not set, the default security zone will be created (if it doesn't exist) and linked. A change of this argument necessitates the re-creation of the resource.

* `gsk_version` - (Optional) The gridscale's Kubernetes version of this instance (e.g. "1.30.3-gs0"). Define which gridscale k8s version will be used to create the cluster. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available gridscale k8s version. **NOTE**: Either `gsk_version` or `release` is set at a time.

* `release` - (Optional) The Kubernetes release of this instance. Define which release will be used to create the cluster. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available releases. **NOTE**: Either `gsk_version` or `release` is set at a time.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `node_pool` - (Optional) The collection of node pool specifications. **NOTE**: Any node pool specification is not yet mutable (except `node_count`).
    * `name` - Name of the node pool.
    * `node_count` - Number of worker nodes.
    * `cores` - Cores per worker node.
    * `memory` - Memory per worker node (in GiB).
    * `storage` - Storage per worker node (in GiB).
    * `storage_type` - Storage type (one of storage, storage_high, storage_insane).
    * `rocket_storage` - Rocket storage per worker node (in GiB).
* `surge_node` - Enable surge node to avoid resources shortage during the cluster upgrade (Default: true).
* `cluster_cidr` - (Immutable) The cluster CIDR that will be used to generate the CIDR of nodes, services, and pods. The allowed CIDR prefix length is /16. If the cluster CIDR is not set, the cluster will use "10.244.0.0/16" as it default (even though the `cluster_cidr` in the k8s resource is empty).
* `cluster_traffic_encryption` - Enables cluster encryption via wireguard if true. Only available for GSK version 1.29 and above. Default is false.

* `oidc_enabled` - (Optional) Enable OIDC for the k8s cluster.

* `oidc_issuer_url` - (Optional) URL of the provider that allows the API server to discover public signing keys. Only URLs that use the https:// scheme are accepted.

* `oidc_client_id` - (Optional) A client ID that all tokens must be issued for.

* `oidc_username_claim` - (Optional) JWT claim to use as the user name.

* `oidc_groups_claim` - (Optional) JWT claim to use as the user's group.

* `oidc_signing_algs` - (Optional)The signing algorithms accepted. Default is 'RS256'. Other option is 'RS512'.

* `oidc_groups_prefix` - (Optional) Prefix prepended to group claims to prevent clashes with existing names (such as system: groups). For example, the value oidc: will create group names like oidc:engineering and oidc:infra.

* `oidc_username_prefix` - (Optional) Prefix prepended to username claims to prevent clashes with existing names (such as system: users). For example, the value oidc: will create usernames like oidc:jane.doe. If this flag isn't provided and --oidc-username-claim is a value other than email the prefix defaults to ( Issuer URL )# where ( Issuer URL ) is the value of --oidc-issuer-url. The value - can be used to disable all prefixing.

* `oidc_required_claim` - (Optional) A key=value pair that describes a required claim in the ID Token. Multiple claims can be set like this: key1=value1,key2=value2.

* `oidc_ca_pem` - (Optional) Custom CA from customer in pem format as string.

* `k8s_hubble` - (Optional) Enable Hubble for the k8s cluster.


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
* `cluster_traffic_encryption` - See Argument Reference above.
* `oidc_enabled` - See Argument Reference above.
* `oidc_issuer_url` - See Argument Reference above.
* `oidc_client_id` - See Argument Reference above.
* `oidc_username_claim` - See Argument Reference above.
* `oidc_groups_claim` - See Argument Reference above.
* `oidc_signing_algs` - See Argument Reference above.
* `oidc_groups_prefix` - See Argument Reference above.
* `oidc_username_prefix` - See Argument Reference above.
* `oidc_required_claim` - See Argument Reference above.
* `oidc_ca_pem` - See Argument Reference above.
* `k8s_hubble` - See Argument Reference above.
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `status` - status indicates the status of the object.
