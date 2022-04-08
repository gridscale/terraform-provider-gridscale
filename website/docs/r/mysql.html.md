---
layout: "gridscale"
page_title: "gridscale: gridscale_mysql"
sidebar_current: "docs-gridscale-resource-mysql"
description: |-
  Manage a MySQL service in gridscale.
---

# gridscale_mysql

Provides a MySQL resource. This can be used to create, modify, and delete MySQL instances.

## Example

The following example shows how one might use this resource to add a MySQL service to gridscale:

```terraform
resource "gridscale_mysql" "terra-mysql-test" {
  name = "my mysql"
	release = "5.7"
	performance_class = "insane"
  max_core_count = 20
  mysql_query_cache_limit = "2M"
	mysql_default_time_zone = "Europe/Berlin"
	mysql_server_id = 2
	mysql_binlog_format = "STATEMENT"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release` - (Required) The mysql release of this instance. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available mysql service releases.

* `performance_class` - (Required) Performance class of mysql service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `mysql_log_bin` - (Optional) mysql parameter: Binary Logging. Default: false.

* `mysql_sql_mode` - (Optional) mysql parameter: SQL Mode. Default: "ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION".

* `mysql_server_id` - (Optional) mysql parameter: Server Id. Default: 1.

* `mysql_query_cache` - (Optional) mysql parameter: Enable query cache. Default: true.

* `mysql_binlog_format` - (Optional) mysql parameter: Binary Logging Format. Default: "ROW".

* `mysql_max_connections` - (Optional) mysql parameter: Max Connections. Default: 4000.

* `mysql_query_cache_size` - (Optional) mysql parameter: Query Cache Size. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)). Default: 128M.

* `mysql_default_time_zone` - (Optional) mysql parameter: Server Timezone. Default: UTC.

* `mysql_query_cache_limit` - (Optional) mysql parameter: Query Cache Limit. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)). Default: 1M.

* `mysql_max_allowed_packet` - (Optional) mysql parameter: Max Allowed Packet Size. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)). Default: 64M.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `network_uuid` - (Optional) The UUID of the network that the service is attached to.

* `security_zone_uuid` -  *DEPRECATED* (Optional, Forcenew) The UUID of the security zone that the service is attached to.

* `max_core_count` - (Optional) Maximum CPU core count. The mysql instance's CPU core count will be autoscaled based on the workload. The number of cores stays between 1 and `max_core_count`.

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
* `mysql_log_bin` - See Argument Reference above.
* `mysql_sql_mode` - See Argument Reference above.
* `mysql_server_id` - See Argument Reference above.
* `mysql_query_cache` - See Argument Reference above.
* `mysql_binlog_format` - See Argument Reference above.
* `mysql_max_connections` - See Argument Reference above.
* `mysql_query_cache_size` - See Argument Reference above.
* `mysql_default_time_zone` - See Argument Reference above.
* `mysql_query_cache_limit` - See Argument Reference above.
* `mysql_max_allowed_packet` - See Argument Reference above.
* `username` - Username for PaaS service. It is used to connect to the mysql instance.
* `password` - Password for PaaS service. It is used to connect to the mysql instance.
* `listen_port` - The port numbers where this mysql service accepts connections.
  * `name` - Name of a port.
  * `host` - Host address.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` -  The UUID of the network that the service is attached to or network UUID containing security zone.
* `service_template_uuid` - PaaS service template that mysql service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Date time this service has been created.
* `status` - Current status of PaaS service.
* `max_core_count` - See Argument Reference above.
* `labels` - See Argument Reference above.
