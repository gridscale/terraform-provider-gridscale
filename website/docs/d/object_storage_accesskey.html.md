---
layout: "gridscale"
page_title: "gridscale: object storage access key"
sidebar_current: "docs-gridscale-datasource-object-storage-accesskey"
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
* `comment` - A comment for the object storage access key.
* `user_uuid` - If a `user_uuid` is set, a user-specific key will get created. If no `user_uuid` is set along a user with write-access to the contract will still only create a user-specific key for themselves while a user with admin-access to the contract will create a contract-level admin key.
* `access_key` - Access key of an object storage.
* `secret_key` - Secret key of an object storage.
