---
layout: "gridscale"
page_title: "gridscale: template"
sidebar_current: "docs-gridscale-datasource-template"
description: |-
  Gets data of a template by name.
---

# gridscale_template

Get data of a template with a specific name. This can be used to make it more visible which template is being used for new storages.

An error is triggered if the template name does not exist.

## Example Usage

Get the template:

```terraform
   data "gridscale_template" "ubuntu" {
     name = "Ubuntu 18.04 LTS"
   }
```

Using the template datasource for the creation of a storage:

```terraform
resource "gridscale_storage" "storage-test"{
  name = "terra-storage-test"
  capacity = 10
  template {
    sshkeys = [ "e17e8fd2-0797-4a00-a85d-eb9a612a6e4e" ]
    template_uuid = data.gridscale_template.ubuntu.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The exact name of the template as show in [the page Template](https://my.gridscale.io/Template).

## Attributes Reference

The following attributes are exported:

* `name` - The name of the template.
* `id` - The UUID of the template.
* `location_uuid` - The location this object is placed.
* `location_country` - Two digit country code (ISO 3166-2) of the location where this object is placed.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The human-readable name of the location. It supports the full UTF-8 character set, with a maximum of 64 characters.
* `status` - Status indicates the status of the object.
* `ostype` - The operating system installed in the template.
* `version` - The version of the template.
* `private` - The object is private, the value will be true. Otherwise the value will be false.
* `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
* `create_time` - The date and time the object was initially created.
* `change_time` - The date and time of the last object change.
* `distro` - The OS distribution that the template contains.
* `description` - Description of the template.
* `usage_in_minutes` - Total minutes the object has been running.
* `capacity` - The capacity of a storage/ISO Image/template/snapshot in GB.
* `current_price` - Defines the price for the current period since the last bill.
* `labels` - List of labels.
