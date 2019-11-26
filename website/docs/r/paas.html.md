---
layout: "gridscale"
page_title: "gridscale: gridscale_paas"
sidebar_current: "docs-gridscale-resource-paas"
description: |-
  Manages a PaaS in gridscale.
---

# gridscale_paas

Provides a PaaS resource. This can be used to create, modify and delete PaaS.

## Example

The following example shows how one might use this resource to add a PaaS to gridscale:

```terraform
resource "gridscale_paas" "terra-paas-test" {
  name = "terra-paas-test"
  service_template_uuid = "f9625726-5ca8-4d5c-b9bd-3257e1e2211a"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `service_template_uuid` - (Required) The template used to create the service.

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `security_zone_uuid` - (Optional) The UUID of the security zone that the service is running in.

* `parameters` - (Optional) Contains the service parameters for the service.
    
    * `param` - (Required) Name of parameter.

    * `value` - (Required) Value of the corresponding parameter.
    
* `resource_limit` - (Optional) A list of service resource limits..
    
    * `resource` - (Required) The name of the resource you would like to cap.

    * `limit` - (Required) The maximum number of the specific resource your service can use.

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
