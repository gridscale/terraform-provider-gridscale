---
layout: "gridscale"
page_title: "gridscale: template"
sidebar_current: "docs-gridscale-datasource-template"
description: |-
  Gets the id of a template by name.
---

# gridscale_template

Get the id of a template with a specific name. This can be used to make it more visible which template is being used for new storages.

An error is triggered if the template name does not exist.

## Example Usage

Get the template:

```hcl
   data "gridscale_template" "ubuntu" {
     name = "Ubuntu 18.04 LTS"
   }
```

Using the template datasource for the creation of a storage:

```hcl
resource "gridscale_storage" "storage-test"{
	name = "terra-storage-test"
	capacity = 10
	template {
		sshkeys = [ "e17e8fd2-0797-4a00-a85d-eb9a612a6e4e" ]
		template_uuid = "${data.gridscale_template.ubuntu.id}"
	}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The exact name of the template as show in [the expert panel of gridscale](https://my.gridscale.io/Expert/Template).

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the template.
