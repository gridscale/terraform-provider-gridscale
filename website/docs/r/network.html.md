---
layout: "gridscale"
page_title: "gridscale: gridscale_network"
sidebar_current: "docs-gridscale-resource-network"
description: |-
  Manages a network in gridscale.
---

# gridscale_network

Provides a network resource. This can be used to create, modify and delete networks.

## Example Usage

The following example shows how one might use this resource to add a network to gridscale:

```terraform
resource "gridscale_network" "networkname"{
  project = "default"
	name = "terraform-network"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `l2security` - (Optional) Defines information about MAC spoofing protection (filters layer2 and ARP traffic based on MAC source). It can only be (de-)activated on a private network - the public network always has l2security enabled. It will be true if the network is public, and false if the network is private.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].


## Attributes

This resource exports the following attributes:

* `project` - The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.
* `id` - The UUID of the network.
* `name` - See Argument Reference above.
* `location_uuid` - Helps to identify which datacenter an object belongs to. The location of the resource depends on the location of the project.
* `l2security` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `network_type` - The type of this network, can be mpls, breakout or network.
* `location_country` - Formatted by the 2 digit country code (ISO 3166-2) of the host country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `public_net` - Is the network public or not.
* `delete_block` - If deleting this network is allowed.
