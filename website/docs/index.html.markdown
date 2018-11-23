---
layout: "gridscale"
page_title: "Provider: gridscale"
sidebar_current: "docs-gridscale-index"
description: |-
  The gridscale provider is used to interact with many resources supported by gridscale.
---

# Gridscale Provider

The gridscale provider is used to interact with many resources supported by gridscale. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available data sources.

## Example Usage

```hcl
# Configure the gridscale provider
provider "gridscale" {
	uuid = "User-UUID"
	token = "API-Token"
}

# Create a server
resource "gridscale_server" "servername"{
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) This is the User-UUID for the gridscale API. It can be found [in the panel](https://my.gridscale.io/APIs/).
* `token` - (Required) This is an API-Token for the gridscale API. It can be created [in the panel](https://my.gridscale.io/APIs/). The created token needs to have full access to be usable by Terraform.

