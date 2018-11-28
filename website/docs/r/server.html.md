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

```hcl
resource "gridscale_server" "terra-server-test"{
	name = "terra-server-test"
	cores = 2
	memory = 1
	storages = [
		"${gridscale_storage.terra-storage-test.id}",
		"UUID of storage 2"
	]
	networks = [
		"${gridscale_network.terra-network-test.id}"
		"UUID of network 2"
	]
	ipv4 = "${gridscale_ipv4.terra-ipv4-test.id}"
	ipv6 = "UUID of ipv6 address"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `cores` - (Required) The number of server cores.

* `memory` - (Required) The amount of server memory in GB.

* `location_uuid` - (Optional) Helps to identify which datacenter an object belongs to. Frankfurt is the default.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `hardware_profile` - (Optional) The hardware profile of the Server. Options are default, legacy, nested, cisco_csr, sophos_utm, f5_bigip and q35 at the moment of writing. Check the 

* `storages` - (Optional) List of the UUIDs of the storages to connect to the server in the format [ "UUID of storage 1" , "UUID of storage 2" ]. The first storage in the list will be set as boot device.

* `networks` - (Optional) List of the UUIDs of the networks to connect to the server in the format [ "UUID of network 1" , "UUID of network 2" ]. The first network in the list will be set as boot device.

* `ipv4` - (Optional) The UUID of the IPv4 address of the server. When this option is set, the server will automatically be connected to the public network, giving it access to the internet.

* `ipv6` - (Optional) The UUID of the IPv6 address of the server. When this option is set, the server will automatically be connected to the public network, giving it access to the internet.

* `power` - (Optional) The power state of the server. Set this to true to will boot the server, false will shut it down.

* `availability_zone` - (Optional) Defines which Availability-Zone the Server is placed.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `cores` - See Argument Reference above.
* `memory` - See Argument Reference above.
* `location_uuid` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `hardware_profile` - See Argument Reference above.
* `storages` - See Argument Reference above.
* `networks` - See Argument Reference above.
* `ipv4` - See Argument Reference above.
* `ipv6` - See Argument Reference above.
* `power` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `location_country` - Formatted by the 2 digit country code (ISO 3166-2) of the host country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `current_price` - The price for the current period since the last bill.
