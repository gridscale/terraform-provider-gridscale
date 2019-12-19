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
}

data "gridscale_object_storage_accesskey" "foo" {
	resource_id   = "${gridscale_object_storage_accesskey.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) ID of a resource (access key of an object storage).

## Attributes Reference

The following attributes are exported:

* `id` - The access key of the object storage.
* `access_key` - Access key of an object storage.
* `secret_key` - Secret key of an object storage.
