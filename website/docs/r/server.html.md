---
layout: "gridscale"
page_title: "gridscale: gridscale_server"
sidebar_current: "docs-gridscale-resource-server"
description: |-
  Manages a server in gridscale.
---

# gridscale_server

Provides a server resource. This can be used to create, modify and delete servers.

## Example Usage

The following example shows how one might use this resource to add a server to gridscale:

```hcl
resource "gridscale_server" "terra-server-test"{
	name = "terra-server-test"
	cores = 2
	memory = 1
	storage {
		object_uuid = "${gridscale_storage.terra-storage-test.id}"
		bootdevice = true
	}
	storage {
    		object_uuid = "UUID of storage 2",
    	}
	network {
		object_uuid = "${gridscale_network.terra-network-test.id}"
		bootdevice = true
	}
	network {
    		object_uuid = "UUID of network 2"
    }
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

* `ipv4` - (Optional) The UUID of the IPv4 address of the server. When this option is set, the server will automatically be connected to the public network, giving it access to the internet.

* `ipv6` - (Optional) The UUID of the IPv6 address of the server. When this option is set, the server will automatically be connected to the public network, giving it access to the internet.

* `power` - (Optional) The power state of the server. Set this to true to will boot the server, false will shut it down.

* `availability_zone` - (Optional) Defines which Availability-Zone the Server is placed.

* `storage` - (Optional) Connects a storage to the server.

    * `object_uuid` - (Required) The object UUID or id of the storage.
    
    * `bootdevice` - (Optional) Make this storage the boot device. This can only be set for one storage.

* `storage` - (Optional) Connects a storage to the server.

    * `object_uuid` - (Required) The object UUID or id of the network.
    
    * `bootdevice` - (Optional) Make this network the boot device. This can only be set for one network.

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
