---
layout: "gridscale"
page_title: "gridscale: gridscale_redis_cache"
sidebar_current: "docs-gridscale-resource-redis-cache"
description: |-
  Manage a Redis cache service in gridscale.
---

# gridscale_redis_cache

Provides a Redis cache resource. This can be used to create, modify, and delete Redis cache instances.

## Example

The following example shows how one might use this resource to add a Redis cache service to gridscale:

```terraform
resource "gridscale_redis_cache" "terra-redis-cache-test" {
  name = "test"
  release = "5.0"
  performance_class = "standard"
  labels = ["test"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `release` - (Required) The Redis cache release of this instance. For convenience, please use [gscloud](https://github.com/gridscale/gscloud) to get the list of available Redis cache service releases.

* `performance_class` - (Required) Performance class of Redis cache service. Available performance classes at the time of writing: `standard`, `high`, `insane`, `ultra`.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `security_zone_uuid` - (Optional) The UUID of the security zone that the service is attached to.

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
* `username` - Username for PaaS service. It is used to connect to the Redis cache instance.
* `password` - Password for PaaS service. It is used to connect to the Redis cache instance.
* `listen_port` - The port numbers where this Redis cache service accepts connections.
  * `name` - Name of a port.
  * `host` - Host address.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - PaaS service template that Redis cache service uses.
* `usage_in_minutes` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Date time this service has been created.
* `status` - Current status of PaaS service.
* `labels` - See Argument Reference above.
