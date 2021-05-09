---
layout: "gridscale"
page_title: "gridscale: imported marketplace application"
sidebar_current: "docs-gridscale-resource-marketplace-application-import"
description: |-
  Manages imported marketplace applications in Gridscale.
---

# gridscale_marketplace_application_import

Provides an imported marketplace application resource. This can be used to import marketplace applications, and delete imported marketplace applications.

## Example Usage

```terraform
resource "gridscale_marketplace_application" "foo" {
  name = "example"
  object_storage_path = "s3://testsnapshot/test.gz"
  category = "Archiving"
  setup_cores = 1
  setup_memory = 1
  setup_storage_capacity = 1
}

resource "gridscale_marketplace_application_import" "importedFoo" {
  import_unique_hash = gridscale_marketplace_application.foo.unique_hash
}
```

## Argument Reference

The following arguments are supported:

* `import_unique_hash` - (Required) Hash of a specific marketplace application that you want to import. A change of this argument necessitates the re-creation of the resource.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the marketplace application.
* `name` - The name of the marketplace application.
* `category` - Category of marketplace application.
* `object_storage_path` - Path to the images for the application.
* `setup_cores` - Number of server's cores.
* `setup_memory` - The capacity of server's memory in GB.
* `setup_storage_capacity` - The capacity of server's storage in GB.
* `meta_license` - License number.
* `meta_os` - Operating system.
* `meta_components` - Components (e.g: MySql, Apache, etc.).
* `meta_overview` - Describes the main function of the application.
* `meta_hints` - Hints.
* `meta_terms_of_use` - Terms of use.
* `meta_icon` - base64 encoded image of the icon.
* `meta_features` - List of functions.
* `meta_author` - Author.
* `meta_advices` - User manual; Wiki URL; ...
* `unique_hash` - Unique hash to allow user to import the self-created marketplace application.
* `is_application_owner` - Whether the you are the owner of application or not.
* `is_published` - Whether the template is published by the partner to their tenant".
* `published_date` - The date when the template is published into other tenant in the same partner.
* `is_publish_requested` - Whether the tenants want their template to be published or not.
* `publish_requested_date` - The date when the tenant requested their template to be published.
* `is_publish_global_requested` - Whether a partner wants their tenant template published to other partners.
* `publish_global_requested_date` - The date when a partner requested their tenants template to be published.
* `is_publish_global` - Whether a template is published to other partner or not.
* `published_global_date` - The date when a template is published to other partner.
* `type` - The type of template.
* `status` - The status of the marketplace application.
* `create_time` - The date and time the marketplace application was initially created.
* `change_time` - The date and time of the last marketplace application change.
