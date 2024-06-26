---
layout: "gridscale"
page_title: "gridscale: gridscale_postgresql"
sidebar_current: "docs-gridscale-resource-postgresql"
description: |-
  Manage a PostgreSQL service in gridscale.
---

# gridscale_postgresql

Provides a PostgreSQL resource. This can be used to create, modify, and delete PostgreSQL instances.

## Example

The following example shows how one might use this resource to add a PostgreSQL service to gridscale:

```terraform
resource "gridscale_postgresql" "terra-postgres-test" {
  name = "test"
  release = "13"
  performance_class = "standard"
  max_core_count = 20
  labels = ["test"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release` - (Required) The PostgreSQL release of this instance. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available PostgreSQL service releases.

* `performance_class` - (Required) Performance class of PostgreSQL service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `network_uuid` - (Optional) The UUID of the network that the service is attached to.

* `security_zone_uuid` -  *DEPRECATED* (Optional, Forcenew) The UUID of the security zone that the service is attached to.

* `max_core_count` - (Optional) Maximum CPU core count. The PostgreSQL instance's CPU core count will be autoscaled based on the workload. The number of cores stays between 1 and `max_core_count`.

* `pgaudit_log_bucket` - (Optional) Object Storage bucket to upload audit logs to. For pgAudit to be enabled these additional parameters need to be configured: pgaudit_log_server_url, pgaudit_log_access_key, pgaudit_log_secret_key.

* `pgaudit_log_server_url` - (Optional) Object Storage server URL the bucket is located on.

* `pgaudit_log_access_key` - (Optional) Access key used to authenticate against Object Storage server.

* `pgaudit_log_secret_key` - (Optional) Secret key used to authenticate against Object Storage server.

* `pgaudit_log_rotation_frequency` - (Optional) Rotation (in minutes) for audit logs. Logs are uploaded to Object Storage once rotated. Default is 5 minutes.

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
* `username` - Username for PaaS service. It is used to connect to the PostgreSQL instance.
* `password` - Password for PaaS service. It is used to connect to the PostgreSQL instance.
* `listen_port` - The port numbers where this PostgreSQL service accepts connections.
  * `name` - Name of a port.
  * `host` - Host address.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` -  The UUID of the network that the service is attached to or network UUID containing security zone.
* `service_template_uuid` - PaaS service template that PostgreSQL service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Date time this service has been created.
* `status` - Current status of PaaS service.
* `max_core_count` - See Argument Reference above.
* `pgaudit_log_bucket` - See Argument Reference above.
* `pgaudit_log_server_url` - See Argument Reference above.
* `pgaudit_log_access_key` - See Argument Reference above.
* `pgaudit_log_secret_key` - See Argument Reference above.
* `pgaudit_log_rotation_frequency` - See Argument Reference above.
* `labels` - See Argument Reference above.
