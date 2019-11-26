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

* `name` - See Argument Reference above.
* `username` - Username for PaaS service.
* `password` - Password for PaaS service.
* `listen_port` - Ports that PaaS service listens to.
    * `name` - Name of a port.
    * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - See Argument Reference above.
* `usage_in_minute` - Number of minutes that PaaS service is in use.
* `current_price` - Current price of PaaS service.
* `change_time` - Time of the last change.
* `create_time` - Time of the creation.
* `status` - Current status of PaaS service.
* `parameter` - See Argument Reference above.
    * `param` - See Argument Reference above.
    * `value` - See Argument Reference above.
* `resource_limit` - See Argument Reference above.
    * `resource` - See Argument Reference above.
    * `limit` - See Argument Reference above.
* `labels` - See Argument Reference above.
