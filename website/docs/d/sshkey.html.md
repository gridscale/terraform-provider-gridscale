---
layout: "gridscale"
page_title: "gridscale: sshkey"
sidebar_current: "docs-gridscale-datasource-sshkey"
description: |-
  Gets the id of an sshkey.
---

# gridscale_sshkey

Get the id of an sshkey resource. This can be used to link SSH keys to a storage when an official template is used.

## Example Usage

Using the sshkey datasource for the creation of a storage:

```terraform
resource "gridscale_sshkey" "sshkey-john"{
	name = "john's computer"
	sshkey = "an ssh public key"
}

resource "gridscale_sshkey" "sshkey-jane"{
	name = "jane's computer"
	sshkey = "an ssh public key"
}

resource "gridscale_storage" "storagename"{
	name = "terraform-storage"
	capacity = 10
	template {
		sshkeys = [
		    "${gridscale_sshkey.sshkey-john.id}",
		    "${gridscale_sshkey.sshkey-jane.id}"
		]
		template_uuid = "4db64bfc-9fb2-4976-80b5-94ff43b1233a"
	}
}
```

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the sshkey.
