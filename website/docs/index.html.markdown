---
layout: "gridscale"
page_title: "Provider: gridscale"
sidebar_current: "docs-gridscale-index"
description: |-
  The gridscale provider is used to interact with many resources supported by gridscale.
---

# gridscale Provider

The gridscale provider is used to interact with many resources supported by gridscale. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available data sources and resources.

## Example Usage

```terraform
# Configure the gridscale provider
provider "gridscale" {
	uuid = var.gridscale_uuid
	token = var.gridscale_token

}

# Create a server
resource "gridscale_server" "servername"{
  # ...
}
```

Also make sure to check out our other Terraform examples over at [github.com/gridscale/terraform_examples](https://github.com/gridscale/terraform_examples).

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) This is the User-UUID for the gridscale API. It can be found [in the panel](https://my.gridscale.io/APIs/). If omitted, the GRIDSCALE_UUID environment variable is used.
* `token` - (Required) This is an API-Token for the gridscale API. It can be created [in the panel](https://my.gridscale.io/APIs/). The created token needs to have full access to be usable by Terraform. If omitted, the GRIDSCALE_TOKEN environment variable is used.
* `api_url` - (Optional) The URL for the API. By default this is set to "https://api.gridscale.io". Do not add a "/" character at the end.
* `gsc_timeout_secs` - (Optional) The general timeout for gsclient-go in seconds. (The specific timeout of a resource's specific operation can only be set lower than/equal to `gsc_timeout_secs`. If it is larger than `gsc_timeout_secs`, `gsc_timeout_secs` will be the timeout of that operation).
