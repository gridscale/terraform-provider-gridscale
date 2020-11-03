---
layout: "gridscale"
page_title: "gridscale: sshkey"
sidebar_current: "docs-gridscale-datasource-sshkey"
description: |-
  Gets data of an SSH key.
---

# gridscale_sshkey

Get data of an sshkey resource. This can be used to link SSH keys to a storage when an official template is used.

## Example Usage

Using the sshkey datasource for the creation of a storage:

```terraform
data "gridscale_sshkey" "sshkey-john"{
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}

data "gridscale_sshkey" "sshkey-jane"{
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}

resource "gridscale_storage" "storagename"{
  name = "terraform-storage"
  capacity = 10
  template {
    sshkeys = [
      data.gridscale_sshkey.sshkey-john.id,
      data.gridscale_sshkey.sshkey-jane.id
    ]
    template_uuid = "4db64bfc-9fb2-4976-80b5-94ff43b1233a"
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the SSH key.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the SSH key.
* `name` - The human-readable name of the SSH key.
* `sshkey` - The OpenSSH public key string of the SSH key.
* `status` - The status of the SSH key.
* `create_time` - The date and time that this SSH key was initially created.
* `change_time` - The date and time of the last change.
* `labels` - The list of labels.
