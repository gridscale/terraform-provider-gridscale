---
layout: "gridscale"
page_title: "gridscale: storage"
sidebar_current: "docs-gridscale-datasource-storage"
description: |-
  Gets the id of a storage.
---

# gridscale_storage

Get the id of a storage resource. This can be used to link storages to a server.

## Example Usage

Using the storage datasource for the creation of a server:

```terraform
resource "gridscale_storage" "storagename"{
	name = "terraform-storage"
	capacity = 10
}

resource "gridscale_server" "servername"{
	name = "terra-server"
	cores = 2
	memory = 4
	storage {
		object_uuid = "${gridscale_storage.storagename.id}"
		bootdevice = true
	}
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the storage.
