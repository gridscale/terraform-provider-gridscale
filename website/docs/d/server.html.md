---
layout: "gridscale"
page_title: "gridscale: gridscale_server"
sidebar_current: "docs-gridscale-datasource-server"
description: |-
  Gets data of a server by its UUID.
---

# gridscale_server

Get data of a server by its UUID.

## Example

```terraform
resource "gridscale_ipv4" "foo1" {
  name   = "newname"
}
resource "gridscale_network" "foo" {
  name   = "newname"
}
resource "gridscale_storage" "foo1" {
  name   = "newname"
  capacity = 1
}
resource "gridscale_server" "foo" {
  name   = "newname"
  cores = 1
  memory = 1
  power = true
  ipv4 = gridscale_ipv4.foo1.id
  network {
    object_uuid = gridscale_network.foo.id
    rules_v4_in {
      order = 0
      protocol = "tcp"
      action = "drop"
      dst_port = "20:80"
      comment = "test"
    }
    rules_v6_in {
      order = 1
      protocol = "tcp"
      action = "drop"
      dst_port = "10:20"
      comment = "test1"
    }
  }
  storage {
    object_uuid = gridscale_storage.foo1.id
  }
}


data "gridscale_server" "foo" {
  resource_id   = gridscale_server.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the firewall.

## Attributes Reference

This resource exports the following attributes:

* `id` - UUID of the server.
* `name` - The name of the server.
* `cores` - The number of server cores.
* `memory` - The amount of server memory in GB.
* `location_uuid` - The location this server is placed. The location of a resource is determined by it's project.
* `labels` - List of labels in the format [ "label1", "label2" ].
* `hardware_profile` - The hardware profile of the server.
* `storage` - Connects a storage to the server.
    * `object_uuid` - The object UUID or id of the storage.
    * `storage_type` - Indicates the speed of the storage. This may be (storage, storage_high or storage_insane).
    * `bootdevice` - True is the storage is a boot device.
    * `object_name` - Name of the storage.
    * `create_time` - Defines the date and time the object was initially created.
    * `capacity` - Capacity of the storage (GB).
    * `controller` - Defines the SCSI controller id. The SCSI defines transmission routes such as Serial Attached SCSI (SAS), Fibre Channel and iSCSI.
    * `bus` - The SCSI bus id. The SCSI defines transmission routes like Serial Attached SCSI (SAS), Fibre Channel and iSCSI. Each SCSI device is addressed via a specific number. Each SCSI bus can have multiple SCSI devices connected to it.
    * `target` - Defines the SCSI target ID. The target ID is a device (e.g. disk).
    * `lun` - Is the common SCSI abbreviation of the Logical Unit Number. A lun is a unique identifier for a single disk or a composite of disks.
    * `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
    * `last_used_template` - Indicates the UUID of the last used template on this storage (inherited from snapshots).
* `network` - Connects a network to the server.
    * `object_uuid` - The object UUID or id of the network.
    * `bootdevice` - Make this network the boot device. This can only be set for one network.
    * `ip` - DHCP IP which is manually assigned to the server (if applicable).
    * `auto_assigned_ip` - DHCP IP which is automatically assigned to the server (if applicable).
    * `object_name` - Name of the network.
    * `ordering` - *DEPRECATED* Defines the ordering of the network interfaces. Lower numbers have lower PCI-IDs.
    * `create_time` - Defines the date and time the object was initially created.
    * `network_type` - One of network, network_high, network_insane.
    * `mac` - network_mac defines the MAC address of the network interface.
    * `firewall_template_uuid` - The UUID of firewall template.
    * `rules_v4_in` - Firewall template rules for inbound traffic - covers IPv4 addresses.
        * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default (Only for inbound).
        * `action` - This defines what the firewall will do. Either accept or drop.
        * `protocol` - Either 'udp' or 'tcp'.
        * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `dst_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `comment` - Comment.
    * `rules_v4_out` - Firewall template rules for outbound traffic - covers IPv4 addresses.
        * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2.
        * `action` - This defines what the firewall will do. Either accept or drop.
        * `protocol` - Either 'udp' or 'tcp'.
        * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `dst_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `comment` - Comment.
    * `rules_v6_in` - Firewall template rules for inbound traffic - covers IPv6 addresses.
        * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default (Only for inbound).
        * `action` - This defines what the firewall will do. Either accept or drop.
        * `protocol` - Either 'udp' or 'tcp'.
        * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `dst_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `comment` - Comment.
    * `rules_v6_out` - Firewall template rules for outbound traffic - covers IPv6 addresses.
        * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2.
        * `action` - This defines what the firewall will do. Either accept or drop.
        * `protocol` - Either 'udp' or 'tcp'.
        * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
        * `src_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `dst_cidr` - Either an IPv4/6 address or and IP network in CIDR format. If this field is empty then this service has access to all IP addresses.
        * `comment` - Comment.
* `ipv4` - The UUID of the IPv4 address of the server.
* `ipv6` - The UUID of the IPv6 address of the server.
* `isoimage` - The UUID of an ISO image in gridscale.
* `power` - The power state of the server.
* `availability_zone` - Defines which Availability-Zone the Server is placed.
* `auto_recovery` - If the server should be auto-started in case of a failure.
* `console_token` - The token used by the panel to open the websocket VNC connection to the server console.
* `legacy` - Legacy-Hardware emulation instead of virtio hardware. If enabled, hot-plugging cores, memory, storage, network, etc. will not work, but the server will most likely run every x86 compatible operating system. This mode comes with a performance penalty, as emulated hardware does not benefit from the virtio driver infrastructure.
* `status` - Status indicates the status of the object.
* `usage_in_minutes_memory` - Total minutes of memory used.
* `usage_in_minutes_cores` - Total minutes of cores used.
* `create_time` - Defines the date and time the object was initially created.
* `change_time` - Defines the date and time of the last object change.
* `current_price` - The price for the current period since the last bill.
