---
layout: "gridscale"
page_title: "gridscale: gridscale_mssql"
sidebar_current: "docs-gridscale-resource-mssql"
description: |-
  Manage a MS SQL server service in gridscale.
---

# gridscale_mssql

Provides a MS SQL server resource. This can be used to create, modify, and delete MS SQL server instances.

## Example

The following example shows how one might use this resource to add a MS SQL server service to gridscale:

```terraform
resource "gridscale_mssql" "terra-mssql-test" {
  name = "test"
  release = "2019"
  performance_class = "standard"
  labels = ["test"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release` - (Required) The MS SQL server release of this instance. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available MS SQL server service releases.

* `performance_class` - (Required) Performance class of MS SQL server service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `security_zone_uuid` - (Optional) The UUID of the security zone that the service is running in.


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
* `username` - Username for PaaS service. It is used to connect to the MS SQL server instance.
* `password` - Password for PaaS service. It is used to connect to the MS SQL server instance.
* `listen_port` - The port numbers where this MS SQL server service accepts connections.
  * `name` - Name of a port.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - PaaS service template that MS SQL server service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Date time this service has been created.
* `status` - Current status of PaaS service.
* `labels` - See Argument Reference above.
