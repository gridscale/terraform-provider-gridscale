---
layout: "gridscale"
page_title: "gridscale: firewall"
sidebar_current: "docs-gridscale-resource-firewall"
description: |-
  Manages a firewall in gridscale.
---

# gridscale_firewall

Provides a firewall resource. This can be used to create, modify and delete firewalls.

## Example Usage

```terraform
resource "gridscale_firewall" "foo" {
  name   = "example-firewall"
  rules_v4_in {
	order = 0
	protocol = "tcp"
	action = "drop"
	dst_port = "20:80"
	comment = "some comments"
  }
  rules_v6_in {
	order = 0
	protocol = "tcp"
	action = "drop"
	dst_port = "2000:3000"
	comment = "some comments"
  }
}
```


## Argument Reference

The following arguments are supported:

***Note: `Optional*` means there is at least 1 rule in the firewall. Otherwise, an error will be returned.

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `rules_v4_in` - (Optional*) Firewall template rules for inbound traffic - covers ipv4 addresses.

    * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

    * `action` - (Required) This defines what the firewall will do. Either accept or drop.

    * `protocol` - (Required) Either 'udp' or 'tcp'.

    * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `comment` - (Optional) Comment.
        
* `rules_v4_out` - (Optional*) Firewall template rules for outbound traffic - covers ipv4 addresses.

    * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

    * `action` - (Required) This defines what the firewall will do. Either accept or drop.

    * `protocol` - (Required) Either 'udp' or 'tcp'.

    * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `comment` - (Optional) Comment.

* `rules_v6_in` - (Optional*) Firewall template rules for inbound traffic - covers ipv6 addresses.

    * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

    * `action` - (Required) This defines what the firewall will do. Either accept or drop.

    * `protocol` - (Required) Either 'udp' or 'tcp'.

    * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `comment` - (Optional) Comment.

* `rules_v6_out` - (Optional*) Firewall template rules for outbound traffic - covers ipv6 addresses.

    * `order` - (Required) The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.

    * `action` - (Required) This defines what the firewall will do. Either accept or drop.

    * `protocol` - (Required) Either 'udp' or 'tcp'.

    * `dst_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_port` - (Optional) A Number between 1 and 65535, port ranges are separated by a colon for FTP.

    * `src_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `dst_cidr` - (Optional) Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.

    * `comment` - (Optional) Comment. 

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is the value of global `timeout`) Used for Creating resource.
* `update` - (Default value is the value of global `timeout`) Used for Updating resource.
* `delete` - (Default value is the value of global `timeout`) Used for Deleteing resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the firewall.
* `name` - The name of the firewall.
* `rules_v4_in` - Firewall template rules for inbound traffic - covers ipv4 addresses.
    * `order` - The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `rules_v4_out` - Firewall template rules for outbound traffic - covers ipv4 addresses.
    * `order` - The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `rules_v6_in` - Firewall template rules for inbound traffic - covers ipv6 addresses.
    * `order` - The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `rules_v6_out` - Firewall template rules for outbound traffic - covers ipv6 addresses.
    * `order` - The order at which the firewall will compare packets against its rules, a packet will be compared against the first rule, it will either allow it to pass or block it and it won t be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2. Packets that do not match any rules are blocked by default.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `network` - The information about networks which are related to this firewall.
    * `object_uuid` - The object UUID or id of the firewall.
    * `object_name` - Name of the firewall.
    * `network_uuid` - The object UUID or id of the network.
    * `network_name` - Name of the network.
    * `create_time` - The date and time the object was initially created.
* `location_name` - The human-readable name of the location. It supports the full UTF-8 charset, with a maximum of 64 characters.
* `status` - Status indicates the status of the object.
* `private` - The object is private, the value will be true. Otherwise the value will be false.
* `create_time` - The date and time the object was initially created.
* `change_time` - The date and time of the last object change.
* `description` - Description of the firewall.
* `labels` - List of labels.
