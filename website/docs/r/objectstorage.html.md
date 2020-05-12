---
layout: "gridscale"
page_title: "gridscale: object storage"
sidebar_current: "docs-gridscale-resource-object-storage"
description: |-
   Manages an access key of an object storage in gridscale.
---

# gridscale_object_storage_accesskey

Provides an access key resource of an object storage. This can be used to create, modify and delete object storages' access keys. 

## Example Usage

```terraform
resource "gridscale_object_storage_accesskey" "foo" {
   timeouts {
      create="10m"
  }
}
```

## Timeouts

Timeouts configuration options (in seconds):
More info: https://www.terraform.io/docs/configuration/resources.html#operation-timeouts

* `create` - (Default value is "5m" - 5 minutes) Used for Creating resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for Deleteing resource.

## Attributes Reference

The following attributes are exported:

* `id` - The access key of the object storage.
* `access_key` - Access key of an object storage.
* `secret_key` - Secret key of an object storage.
