---
layout: "gridscale"
page_title: "gridscale: gridscale_postgresql"
sidebar_current: "docs-gridscale-resource-postgresql"
description: |-
  Manage a postgreSQL service in gridscale.
---

# gridscale_postgresql

Provides a postgreSQL resource. This can be used to create, modify, and delete postgreSQL instances.

## Example

The following example shows how one might use this resource to add a postgreSQL service to gridscale:

```terraform
resource "gridscale_postgresql" "terra-postgres-test" {
  name = "test"
  release_no = "13"
  performance_class = "standard"
  max_core_count = 20
  labels = ["test"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release_no` - (Required) Release no. of postgreSQL service. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available postgreSQL service release numbers.

* `performance_class` - (Required) Performance class of postgreSQL service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `security_zone_uuid` - (Optional) The UUID of the security zone that the service is running in.

* `max_core_count` - (Optional) Maximum CPU core count. It is used for autoscaling.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "15m" - 15 minutes) Used for creating a resource.
* `update` - (Default value is "15m" - 15 minutes) Used for updating a resource.
* `delete` - (Default value is "15m" - 15 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `username` - Username for PaaS service.
* `password` - Password for PaaS service.
* `listen_port` - Ports that PaaS service listens to.
  * `name` - Name of a port.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - PaaS service template that postgreSQL service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Time of the creation.
* `status` - Current status of PaaS service.
* `max_core_count` - See Argument Reference above.
* `labels` - See Argument Reference above.
