---
layout: "gridscale"
page_title: "gridscale: object storage"
sidebar_current: "docs-gridscale-datasource-object-storage"
description: |-
  Gets the data of an access key of an object storage.
---

# gridscale_object_storage_accesskey

Get data of an access key resource of an object storage.

## Example Usage

```terraform
resource "gridscale_object_storage_accesskey" "foo" {
  project = "default"
}

data "gridscale_object_storage_accesskey" "foo" {
  project = "default"
	resource_id   = "${gridscale_object_storage_accesskey.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.

* `resource_id` - (Required) ID of a resource (access key of an object storage).

## Attributes Reference

The following attributes are exported:

* `project` - The name of project which is set in GRIDSCALE_PROJECTS_TOKENS env variable.
* `id` - The access key of the object storage.
* `access_key` - Access key of an object storage.
* `secret_key` - Secret key of an object storage.
