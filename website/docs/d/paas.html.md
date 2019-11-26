---
layout: "gridscale"
page_title: "gridscale: paas"
sidebar_current: "docs-gridscale-datasource-paas"
description: |-
  Gets the data of a PaaS based on given UUID.
---

# gridscale_paas

Get a PaaS resource based on given UUID.

## Example Usage

Using the network datasource for the creation of a PaaS:

```terraform
resource "gridscale_paas" "foo" {
  name = "foo"
  service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
}

data "gridscale_paas" "foo" {
	resource_id   = "${gridscale_paas.foo.id}"
}
```

## Attributes Reference

The following attributes are exported:

* `name` - The human-readable name of the object.
* `username` - Username for PaaS service.
* `password` - Password for PaaS service.
* `listen_port` - Ports that PaaS service listens to.
    * `name` - Name of a port.
    * `listen_port` - Port number.
* `security_zone_uuid` - The UUID of the security zone that the service is running in.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - The template used to create the service.
* `usage_in_minute` - Number of minutes that PaaS service is in use.
* `current_price` - Current price of PaaS service.
* `change_time` - Time of the last change.
* `create_time` - Time of the creation.
* `status` - Current status of PaaS service.
* `parameter` - Contains the service parameters for the service.
    * `param` - Name of parameter.
    * `value` - Value of the corresponding parameter.
* `resource_limit` - A list of service resource limits.
    * `resource` - The name of the resource you would like to cap.
    * `limit` - The maximum number of the specific resource your service can use.
* `labels` - List of labels in the format [ "label1", "label2" ].
