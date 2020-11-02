---
layout: "gridscale"
page_title: "gridscale: marketplace application"
sidebar_current: "docs-gridscale-resource-marketplace-application"
description: |-
  Manages marketplace applications in Gridscale.
---

# gridscale_marketplace_application

Provides a marketplace application resource. This can be used to create, modify, and delete marketplace applications.

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `object_storage_path` - (Required) Path to the images for the application, must be in .gz format and started with s3//.

* `category` - (Required) Category of marketplace application. Accepted values: "CMS", "project management", "Adminpanel", "Collaboration", "Cloud Storage", "Archiving".

* `setup_cores` - (Required) Number of server's cores.

* `setup_memory` - (Required) The capacity of server's memory in GB.

* `setup_storage_capacity` - (Required) The capacity of server's storage in GB.

* `meta_license` - (Optional) License number.

* `meta_os` - (Optional) Operating system.

* `meta_components` - (Optional) Components (e.g: MySql, Apache, etc.).

* `meta_overview` - (Optional) Describes the main function of the application.

* `meta_hints` - (Optional) Hints.

* `meta_terms_of_use` - (Optional) Terms of use.

* `meta_icon` - (Optional) base64 encoded image of the icon.

* `meta_features` - (Optional) List of functions.

* `meta_author` - (Optional) Author.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the marketplace application.
* `name` - See Argument Reference above.
* `category` - See Argument Reference above.
* `object_storage_path` - See Argument Reference above.
* `setup_cores` - See Argument Reference above.
* `setup_memory` - See Argument Reference above.
* `setup_storage_capacity` - See Argument Reference above.
* `meta_license` - See Argument Reference above.
* `meta_os` - See Argument Reference above.
* `meta_components` - See Argument Reference above.
* `meta_overview` - See Argument Reference above.
* `meta_hints` - See Argument Reference above.
* `meta_terms_of_use` - See Argument Reference above.
* `meta_icon` - See Argument Reference above.
* `meta_features` - See Argument Reference above.
* `meta_author` - See Argument Reference above.
* `meta_advices` - See Argument Reference above.
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
