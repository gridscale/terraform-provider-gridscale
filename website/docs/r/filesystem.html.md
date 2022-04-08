---
layout: "gridscale"
page_title: "gridscale: gridscale_mysql"
sidebar_current: "docs-gridscale-resource-filesystem"
description: |-
  Manage a Filesystem service in gridscale.
---

# gridscale_filesystem

Provides a Filesystem service resource. This can be used to create, modify, and delete Filesystem service instances.

## Example

The following example shows how one might use this resource to add a Filesystem service to gridscale:

```terraform
resource "gridscale_filesystem" "terra-filesystem-test" {
  name = "my filesystem"
  release = "1"
  performance_class = "standard"
  root_squash = true
  allowed_ip_ranges = ["192.14.2.2", "192.168.0.0/16"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release` - (Required) The filesystem service release of this instance. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available filesystem service releases.

* `performance_class` - (Required) Performance class of filesystem service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `network_uuid` - (Optional) The UUID of the network that the service is attached to.

* `security_zone_uuid` -  *DEPRECATED* (Optional, Forcenew) The UUID of the security zone that the service is attached to.

* `root_squash` - (Optional) Map root user/group ownership to anon_uid/anon_gid.

* `allowed_ip_ranges` - (Optional) Allowed CIDR block or IP address in CIDR notation.

* `anon_uid` - (Optional) Target user id when root squash is active.

* `anon_gid` - (Optional) Target group id when root squash is active.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "15m" - 15 minutes) Used for creating a resource.
* `update` - (Default value is "15m" - 15 minutes) Used for updating a resource.
* `delete` - (Default value is "15m" - 15 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `release` - See Argument Reference above.
* `performance_class` - See Argument Reference above.
* `root_squash` - See Argument Reference above.
* `allowed_ip_ranges` - See Argument Reference above.
* `anon_uid` - See Argument Reference above.
* `anon_gid` - See Argument Reference above.
* `listen_port` - The port numbers where the filesystem service accepts connections.
  * `name` - Name of a port.
  * `host` - Host address.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` -  The UUID of the network that the service is attached to or network UUID containing security zone.
* `service_template_uuid` - PaaS service template that filesystem service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Date time this service has been created.
* `status` - Current status of PaaS service.
* `labels` - See Argument Reference above.
