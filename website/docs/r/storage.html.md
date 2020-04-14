---
layout: "gridscale"
page_title: "gridscale: gridscale_storage"
sidebar_current: "docs-gridscale-resource-storage"
description: |-
  Manages a storage in gridscale.
---

# gridscale_storage

Provides a storage resource. This can be used to create, modify and delete storages.

## Example Usage

The following example shows how one might use this resource to add a storage to gridscale:

```terraform
resource "gridscale_storage" "storage-john"{
	name = "john's storage"
	capacity = 10
	storage_type = "storage_high"
	template {
	    template_uuid = "4db64bfc-9fb2-4976-80b5-94ff43b1233a"
	    password = var.gridscale_password-john
	    password_type = "plain"
	    hostname = "Ubuntu"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `capacity` - (Required) required (integer - minimum: 1 - maximum: 4096).

* `storage_type` - (Optional) (one of storage, storage_high, storage_insane).

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

* `template` - (Optional) List of labels in the format [ "label1", "label2" ].

    * `template_uuid` - (Required) The UUID of a template. This can be found in the [expert panel](https://my.gridscale.io/Expert/Template) by clicking more on the template or by using a gridscale_template datasource.

    * `password` - (Optional) The root (Linux) or Administrator (Windows) password to set for the installed storage. Valid only for public templates. The password has to be either plain-text or a crypt string (modular crypt format - MCF).

    * `password_type` - (Optional) (one of plain, crypt) Required if password is set (ignored for private templates and public Windows templates).

    * `sshkeys` - (Optional) (array of any - minItems: 0) Public Linux templates only! The UUIDs of SSH keys to be added for the root user.

    * `hostname` - (Optional) The hostname of the installed server (ignored for private templates and public windows templates).

~> **Note** When using official templates using either a password and password_type or at least one SSH public key is required. This is not the case when using custom templates. For official templates password authentication for SSH is enabled by default, so be sure to pick a strong password.


## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is the value of global `timeout`) Used for Creating resource.
* `update` - (Default value is the value of global `timeout`) Used for Updating resource.
* `delete` - (Default value is the value of global `timeout`) Used for Deleteing resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `capacity` - See Argument Reference above.
* `storage_type` - See Argument Reference above.
* `location_uuid` - Helps to identify which datacenter an object belongs to. The location of the resource depends on the location of the project.
* `labels` - See Argument Reference above.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
* `location_country` - Formatted by the 2 digit country code (ISO 3166-2) of the host country.
* `location_iata` - Uses IATA airport code, which works as a location identifier.
* `location_name` - The location name.
* `license_product_no` - If a template has been used that requires a license key (e.g. Windows Servers) this shows the product_no of the license (see the /prices endpoint for more details).
* `last_used_template` - Indicates the UUID of the last used template on this storage (inherited from snapshots).
* `usage_in_minutes` - The amount of minutes the IP address has been in use.
* `current_price` - The price for the current period since the last bill.
