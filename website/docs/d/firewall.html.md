---
layout: "gridscale"
page_title: "gridscale: firewall"
sidebar_current: "docs-gridscale-datasource-firewall"
description: |-
  Gets data of a firewall by its UUID.
---

# gridscale_firewall

Get data of a firewall by its UUID.

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

data "gridscale_firewall" "foo" {
	resource_id   = gridscale_firewall.foo.id
}
```


## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the firewall.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the firewall.
* `name` - The name of the firewall.
* `rules_v4_in` - Firewall template rules for inbound traffic - covers ipv4 addresses.
    * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `rules_v4_out` - Firewall template rules for outbound traffic - covers ipv4 addresses.
    * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `rules_v6_in` - Firewall template rules for inbound traffic - covers ipv6 addresses.
    * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2.
    * `action` - This defines what the firewall will do. Either accept or drop.
    * `protocol` - Either 'udp' or 'tcp'.
    * `dst_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_port` - A Number between 1 and 65535, port ranges are separated by a colon for FTP.
    * `src_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `dst_cidr` - Either an IPv4/6 address or and IP Network in CIDR format. If this field is empty then this service has access to all IPs.
    * `comment` - Comment.
* `rules_v6_out` - Firewall template rules for outbound traffic - covers ipv6 addresses.
    * `order` - The order at which the firewall will compare packets against its rules. A packet will be compared against the first rule, it will either allow it to pass or block it and it won't be matched against any other rules. However, if it does no match the rule, then it will proceed onto rule 2.
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
