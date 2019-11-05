---
layout: "gridscale"
page_title: "gridscale: gridscale_server"
sidebar_current: "docs-gridscale-resource-server"
description: |-
  Manages a server in gridscale.
---

# gridscale_server

Provides a server resource. This can be used to create, modify and delete servers.

## Example

The following example shows how one might use this resource to add a server to gridscale:

```terraform
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
	isoimage = "9be3e0a3-42ac-4207-8887-3383c405724d"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `cores` - (Required) The number of server cores.

* `memory` - (Required) The amount of server memory in GB.

* `location_uuid` - (Optional, ForceNew) Helps to identify which datacenter an object belongs to. Frankfurt is the default.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `auto_recovery` - (Optional) If the server should be auto-started in case of a failure (default=true).

* `hardware_profile` - (Optional, ForceNew) The hardware profile of the Server. Options are default, legacy, nested, cisco_csr, sophos_utm, f5_bigip and q35 at the moment of writing. Check the

* `ipv4` - (Optional) The UUID of the IPv4 address of the server. When this option is set, the server will automatically be connected to the public network, giving it access to the internet.

* `ipv6` - (Optional) The UUID of the IPv6 address of the server. When this option is set, the server will automatically be connected to the public network, giving it access to the internet.

* `isoimage` - (Optional) The UUID of an ISO image in gridscale. The server will automatically boot from the ISO if one was added. The UUIDs of ISO images can be found in [the expert panel](https://my.gridscale.io/Expert/ISOImage).

* `power` - (Optional, Computed) The power state of the server. Set this to true to will boot the server, false will shut it down.

* `availability_zone` - (Optional, Computed) Defines which Availability-Zone the Server is placed.

* `storage` - (Optional) Connects a storage to the server.

    * `object_uuid` - (Required) The object UUID or id of the storage.
    
    * `storage_type` - (Computed) Indicates the speed of the storage. This may be (storage, storage_high or storage_insane).

    * `bootdevice` - (Computed) True is the storage is a boot device.
    
    * `object_name` - (Computed) Name of the storage.

    * `create_time` - (Computed) Defines the date and time the object was initially created.

    * `capacity` - (Computed) Capacity of the storage (GB).

    * `controller` - (Computed) Defines the SCSI controller id. The SCSI defines transmission routes such as Serial Attached SCSI (SAS), Fibre Channel and iSCSI.

    * `bus` - (Computed) The SCSI bus id. The SCSI defines transmission routes like Serial Attached SCSI (SAS), Fibre Channel and iSCSI. Each SCSI device is addressed via a specific number. Each SCSI bus can have multiple SCSI devices connected to it.
    
    * `target` - (Computed) Defines the SCSI target ID. The target ID is a device (e.g. disk).
    
    * `lun` - (Computed) Is the common SCSI abbreviation of the Logical Unit Number. A lun is a unique identifier for a single disk or a composite of disks.

    * `license_product_no` - (Computed) If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).

    * `last_used_template` - (Computed) Indicates the UUID of the last used template on this storage (inherited from snapshots).

* `network` - (Optional) Connects a network to the server.

    * `object_uuid` - (Required) The object UUID or id of the network.

    * `bootdevice` - (Optional, Computed) Make this network the boot device. This can only be set for one network.

    * `object_name` - (Computed) Name of the network.

    * `ordering` - (Computed) Defines the ordering of the network interfaces. Lower numbers have lower PCI-IDs.

    * `create_time` - (Computed) Defines the date and time the object was initially created.

    * `network_type` - (Computed) One of network, network_high, network_insane.

    * `mac` - (Computed) network_mac defines the MAC address of the network interface.

    * `firewall_template_uuid` - (Optional) The UUID of firewall template.

    * `rules_v4_in` - (Optional) Firewall template rules for inbound traffic - covers ipv4 addresses.
    
        * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

        * `action` - (Required) This defines what the firewall will do. Either accept or drop.

        * `protocol` - (Optional) Either 'udp' or 'tcp'.
        
        * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
        
        * `src_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
        
        * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

        * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

        * `comment` - (Optional) Comment.
        
    * `rules_v4_out` - (Optional) Firewall template rules for outbound traffic - covers ipv4 addresses.
        
            * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.
    
            * `action` - (Required) This defines what the firewall will do. Either accept or drop.
    
            * `protocol` - (Optional) Either 'udp' or 'tcp'.
            
            * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
            
            * `src_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
            
            * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    
            * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    
            * `comment` - (Optional) Comment.

    * `rules_v6_in` - (Optional) Firewall template rules for inbound traffic - covers ipv6 addresses.
    
        * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

        * `action` - (Required) This defines what the firewall will do. Either accept or drop.

        * `protocol` - (Optional) Either 'udp' or 'tcp'.
        
        * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
        
        * `src_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
        
        * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

        * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

        * `comment` - (Optional) Comment.
        
    * `rules_v6_out` - (Optional) Firewall template rules for outbound traffic - covers ipv6 addresses.
    
        * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

        * `action` - (Required) This defines what the firewall will do. Either accept or drop.

        * `protocol` - (Optional) Either 'udp' or 'tcp'.
        
        * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
        
        * `src_port` - (Optional) A Number between 1 and 65535, port ranges are seperated by a colon for FTP.
        
        * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

        * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

        * `comment` - (Optional) Comment.

* `current_price` - (Computed) The price for the current period since the last bill.

* `console_token` - (Computed) The token used by the panel to open the websocket VNC connection to the server console.

* `legacy` - (Computed) Legacy-Hardware emulation instead of virtio hardware. If enabled, hotplugging cores, memory, storage, network, etc. will not work, but the server will most likely run every x86 compatible operating system. This mode comes with a performance penalty, as emulated hardware does not benefit from the virtio driver infrastructure.

* `usage_in_minutes_memory` - (Computed) Total minutes of memory used.

* `usage_in_minutes_cores` - (Computed) Total minutes of cores used.

* `create_time` - (Computed) Defines the date and time the object was initially created.

* `change_time` - (Computed) Defines the date and time of the last object change.

* `status` - (Computed) Status indicates the status of the object.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `cores` - See Argument Reference above.
* `memory` - See Argument Reference above.
* `location_uuid` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `hardware_profile` - See Argument Reference above.
* `storage` - See Argument Reference above.
* `network` - See Argument Reference above.
* `ipv4` - See Argument Reference above.
* `ipv6` - See Argument Reference above.
* `isoimage` - See Argument Reference above.
* `power` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `auto_recovery` - See Argument Reference above.
* `console_token` - See Argument Reference above.
* `legacy` - status indicates the status of the object.
* `status` - status indicates the status of the object.
* `usage_in_minutes_memory` - status indicates the status of the object.
* `usage_in_minutes_cores` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `current_price` - The price for the current period since the last bill.
