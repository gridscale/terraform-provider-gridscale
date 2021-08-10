---
layout: "gridscale"
page_title: "gridscale: gridscale_paas"
sidebar_current: "docs-gridscale-resource-paas"
description: |-
  Manages a PaaS in gridscale.
---

# gridscale_paas

Provides a PaaS resource. This can be used to create, modify, and delete PaaS instances.

## Example

The following example shows how one might use this resource to add a PaaS to gridscale:

```terraform
resource "gridscale_paas" "terra-paas-test" {
  name = "terra-paas-test"
  service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
  timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `service_template_uuid` - (Required) The template used to create the service.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `security_zone_uuid` - (Optional) The UUID of the security zone that the service is running in.

* `parameters` - (Optional) Contains the service parameters for the service.

  * `param` - (Required) Name of parameter.

  * `value` - (Required) Value of the corresponding parameter.

* `resource_limit` - (Optional) A list of service resource limits..

  * `resource` - (Required) The name of the resource you would like to cap.

  * `limit` - (Required) The maximum number of the specific resource your service can use.

  * `type` - (Required) Primitive type of the parameter: bool, int (better use float for int case), float, string.

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
  * `host` - Host address.
  * `listen_port` - Port number.
* `security_zone_uuid` - See Argument Reference above.
* `network_uuid` - Network UUID containing security zone.
* `service_template_uuid` - See Argument Reference above.
* `service_template_uuid_computed` - Template that PaaS service uses. The `service_template_uuid_computed` will be different from `service_template_uuid`, when `service_template_uuid` is updated outside of terraform.
* `usage_in_minute` - Number of minutes that PaaS service is in use.
* `change_time` - Time of the last change.
* `create_time` - Time of the creation.
* `status` - Current status of PaaS service.
* `parameter` - See Argument Reference above.
  * `param` - See Argument Reference above.
  * `value` - See Argument Reference above.
  * `type` - See Argument Reference above.
* `resource_limit` - See Argument Reference above.
  * `resource` - See Argument Reference above.
  * `limit` - See Argument Reference above.
* `labels` - See Argument Reference above.
