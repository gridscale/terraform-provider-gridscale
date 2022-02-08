---
layout: "gridscale"
page_title: "gridscale: location"
sidebar_current: "docs-gridscale-resource-location"
description: |-
  Manages a location in gridscale.
---

# gridscale_location

Provides a location resource. This can be used to create, modify, and delete location.

## Example Usage

The following example shows how one might use this resource to add a location to gridscale:

```terraform
resource "gridscale_location" "foo" {
  name   = "my-location"
  parent_location_uuid = "%s"
  product_no = 1500001
  cpunode_count = 20
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The exact name of the location as show in [the expert panel of gridscale](https://my.gridscale.io/Expert/Template).

* `parent_location_uuid` - (Required, ForceNew) The location_uuid of an existing public location in which to create the private location.

* `cpunode_count` - (Required) The number of dedicated cpunodes to assigne to the private location.

* `product_no` - (Required, ForceNew) The product number of a valid and available dedicated cpunode article.

* `labels` - (Optional) List of labels.

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the location.
* `name` - The name of the location.
* `parent_location_uuid` - See Argument Reference above.
* `cpunode_count` - See Argument Reference above.
* `product_no` - See Argument Reference above.
* `iata` - Uses IATA airport code, which works as a location identifier.
* `country` - The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `labels` - List of labels.
* `active` - True if the location is active.
* `cpunode_count_change_requested` - The requested number of dedicated cpunodes.
* `product_no_change_requested` - The product number of a valid and available dedicated cpunode article.
* `parent_location_uuid_change_requested` - The location_uuid of an existing public location in which to create the private location.
* `public` - True if this location is publicly available or a private location.
* `certification_list` - List of certifications.
* `city` - The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `data_protection_agreement` - Data protection agreement.
* `geo_location` - Geo location.
* `green_energy` - Green energy.
* `operator_certification_list` - List of operator certifications.
* `owner` - The human-readable name of the owner.
* `owner_website` - The website of the owner.
* `site_name` - The human-readable name of the website.
* `hardware_profiles` - List of supported hardware profiles.
* `has_rocket_storage` - TRUE if the location supports rocket storage.
* `has_server_provisioning` - TRUE if the location supports server provisioning.
* `object_storage_region` - The region of the object storage.
* `backup_center_location_uuid` - The location_uuid of a backup location.
