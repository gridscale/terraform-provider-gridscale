---
layout: "gridscale"
page_title: "gridscale: object storage bucket"
sidebar_current: "docs-gridscale-resource-object-storage-bucket"
description: |-
   Manages an object storage bucket in gridscale.
---

# gridscale_object_storage_bucket

Provides an object storage bucket in gridscale. This can be used to create, and delete object storage buckets.

## Example Usage

```terraform
resource "gridscale_object_storage_accesskey" "foo" {
   timeouts {
      create="10m"
  }
}

resource "gridscale_object_storage_bucket" "foo-bucket" {
   access_key = gridscale_object_storage_accesskey.foo.access_key
   secret_key = gridscale_object_storage_accesskey.foo.secret_key
   bucket_name = "my-bucket"
}
```

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Argument Reference

The following arguments are supported:

* `access_key` - (Required, Force New) Access key.
* `secret_key` - (Required, Force New) Secret key.
* `s3_host` - (Required, Force New) Host of the s3. Default: "gos3.io".
* `bucket_name` - (Required, Force New) Name of the bucket.
