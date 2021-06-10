---
layout: "gridscale"
page_title: "gridscale: gridscale_mariadb"
sidebar_current: "docs-gridscale-resource-mariadb"
description: |-
  Manage a MariaDB service in gridscale.
---

# gridscale_mariadb

Provides a MariaDB resource. This can be used to create, modify, and delete MariaDB instances.

## Example

The following example shows how one might use this resource to add a MariaDB service to gridscale:

```terraform
resource "gridscale_mariadb" "terra-mariadb-test" {
  name = "my mariadb"
	release = "10.5"
	performance_class = "insane"
  max_core_count = 20
  mariadb_query_cache_limit = "2M"
	mariadb_default_time_zone = "Europe/Berlin"
	mariadb_server_id = 2
	mariadb_binlog_format = "STATEMENT"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release` - (Required) The MariaDB release of this instance. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available MariaDB service releases.

* `performance_class` - (Required) Performance class of MariaDB service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `mariadb_log_bin` - (Optional) MariaDB parameter: Binary Logging. Default: false.

* `mariadb_sql_mode` - (Optional) MariaDB parameter: SQL Mode. Default: "NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO".

* `mariadb_server_id` - (Optional) MariaDB parameter: Server Id. Default: 1.

* `mariadb_query_cache` - (Optional) MariaDB parameter: Enable query cache. Default: true.

* `mariadb_binlog_format` - (Optional) MariaDB parameter: Binary Logging Format. Default: "MIXED".

* `mariadb_max_connections` - (Optional) MariaDB parameter: Max Connections. Default: 4000.

* `mariadb_query_cache_size` - (Optional) MariaDB parameter: Query Cache Size. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)). Default: 128M.

* `mariadb_default_time_zone` - (Optional) MariaDB parameter: Server Timezone. Default: UTC.

* `mariadb_query_cache_limit` - (Optional) MariaDB parameter: Query Cache Limit. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)). Default: 1M.

* `mariadb_max_allowed_packet` - (Optional) MariaDB parameter: Max Allowed Packet Size. Format: xM (where x is an integer, M stands for unit: k(kB), M(MB), G(GB)). Default: 64M.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `security_zone_uuid` - (Optional) The UUID of the security zone that the service is running in.

* `max_core_count` - (Optional) Maximum CPU core count. The MariaDB instance's CPU core count will be autoscaled based on the workload. The number of cores stays between 1 and `max_core_count`.

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
* `mariadb_log_bin` - See Argument Reference above.
* `mariadb_sql_mode` - See Argument Reference above.
* `mariadb_server_id` - See Argument Reference above.
* `mariadb_query_cache` - See Argument Reference above.
* `mariadb_binlog_format` - See Argument Reference above.
* `mariadb_max_connections` - See Argument Reference above.
* `mariadb_query_cache_size` - See Argument Reference above.
* `mariadb_default_time_zone` - See Argument Reference above.
* `mariadb_query_cache_limit` - See Argument Reference above.
* `mariadb_max_allowed_packet` - See Argument Reference above.
* `username` - Username for PaaS service. It is used to connect to the MariaDB instance.
* `password` - Password for PaaS service. It is used to connect to the MariaDB instance.
* `listen_port` - The port numbers where this MariaDB service accepts connections.
  * `name` - Name of a port.
  * `host` - Host address.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - PaaS service template that MariaDB service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Date time this service has been created.
* `status` - Current status of PaaS service.
* `max_core_count` - See Argument Reference above.
* `labels` - See Argument Reference above.
