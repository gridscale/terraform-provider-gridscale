---
layout: "gridscale"
page_title: "gridscale: location"
sidebar_current: "docs-gridscale-datasource-location"
description: |-
  Get the data of a location in gridscale.
---

# gridscale_location

  Get the data of a location in gridscale.

## Example Usage

```terraform
data "gridscale_location" "foo" {
  resource_id = "45ed677b-3702-4b36-be2a-a2eab9827950"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the location.


## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the location.
* `name` - The name of the location.
* `parent_location_uuid` - The location_uuid of an existing public location in which to create the private location.
* `cpunode_count` - The number of dedicated cpunodes to assigne to the private location.
* `product_no` - The product number of a valid and available dedicated cpunode article.
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
