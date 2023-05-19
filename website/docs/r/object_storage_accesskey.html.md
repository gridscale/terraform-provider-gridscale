---
layout: "gridscale"
page_title: "gridscale: object storage access key"
sidebar_current: "docs-gridscale-resource-object-storage-accesskey"
description: |-
   Manages an access key of an object storage in gridscale.
---

# gridscale_object_storage_accesskey

Provides an access key resource of an object storage. This can be used to create, modify, and delete object storages' access keys.

## Example Usage

```terraform
resource "gridscale_object_storage_accesskey" "foo" {
   comment = "foo test"
   timeouts {
      create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) A comment for the object storage access key.
* `user_uuid` - (Optional, Computed) If a `user_uuid` is set, a user-specific key will get created. If no `user_uuid` is set along a user with write-access to the contract will still only create a user-specific key for themselves while a user with admin-access to the contract will create a contract-level admin key.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes Reference

The following attributes are exported:

* `id` - The access key of the object storage.
* `comment` - See Argument Reference above.
* `user_uuid` - See Argument Reference above.
* `access_key` - Access key of an object storage.
* `secret_key` - Secret key of an object storage.
