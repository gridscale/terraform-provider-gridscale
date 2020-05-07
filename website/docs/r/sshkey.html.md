---
layout: "gridscale"
page_title: "gridscale: gridscale_sshkey"
sidebar_current: "docs-gridscale-resource-sshkey"
description: |-
  Manages an SSH public key in gridscale.
---

# gridscale_sshkey

Provides an SSH public key resource. This can be used to create, modify and delete SSH public keys.

## Example Usage

The following example shows how one might use this resource to add an SSH public key to gridscale:

```terraform
resource "gridscale_sshkey" "sshkey-john"{
	name = "john's computer"
	sshkey = "an ssh public key"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The human-readable name of the object. It supports the full UTF-8 charset, with a maximum of 64 characters.

* `sshkey` - (Required) This is the OpenSSH public key string (all key types are supported => ed25519, ecdsa, dsa, rsa, rsa1).

* `labels` - (Optional) List of labels in the format [ "label1", "label2" ].

## Timeouts

Timeouts configuration options (in seconds):

* `create` - (Default value is 0, it means the operation runs without a timeout) Used for Creating resource.
* `update` - (Default value is 0, it means the operation runs without a timeout) Used for Updating resource.
* `delete` - (Default value is 0, it means the operation runs without a timeout) Used for Deleteing resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `sshkey` - See Argument Reference above.
* `labels` - See Argument Reference above.
* `status` - status indicates the status of the object.
* `create_time` - The time the object was created.
* `change_time` - Defines the date and time of the last object change.
