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

# Declare required provider
terraform {
  required_providers {
    gridscale = {
      source = "gridscale/gridscale"
      version = "1.7.3"
    }
  }
}


# Configure the gridscale provider
provider "gridscale" {
  uuid = var.gridscale_uuid
  token = var.gridscale_token
}

# Create a server
resource "gridscale_server" "servername"{
  # …
}
```

Also make sure to check out our other Terraform examples over at [github.com/gridscale/terraform_examples](https://github.com/gridscale/terraform_examples).

~> **Note** To verify that you are using the official released version of gridscale terraform provider, you just need to check the format of the `source` whether it is `<NAMESPACE>/<NAME>/<PROVIDER>`. For more information, visit [registry.terraform.io/providers/gridscale/gridscale/](https://registry.terraform.io/providers/gridscale/gridscale/latest).

~> **Good to know** To protect any resources from being accidently destroyed, take a look at [lifecycle.prevent_destroy](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html#prevent_destroy). Other arguments in [lifecycle](https://www.terraform.io/docs/language/meta-arguments/lifecycle.html) are also good to know.

## Argument Reference

The following arguments are supported:

* `source` - (Required) The global source address for the provider you intend to use. In this case, we use `gridscale/gridscale`.
* `version` - (Optional) A version constraint specifying which subset of available provider versions the module is compatible with.
* `uuid` - (Optional) This is the User-UUID for the gridscale API. It can be found [in the panel](https://my.gridscale.io/APIs/). If omitted, the GRIDSCALE_UUID environment variable is used.
* `token` - (Optional) This is an API-Token for the gridscale API. It can be created [in the panel](https://my.gridscale.io/APIs/). The created token needs to have full access to be usable by Terraform. If omitted, the GRIDSCALE_TOKEN environment variable is used.
* `api_url` - (Optional) The URL for the API. If omitted, the GRIDSCALE_URL environment variable may be used. By default this is set to "https://api.gridscale.io". Do not add a "/" character at the end.
* `http_headers` - (Optional) Custom HTTP headers sent to gridscale API. If omitted, the GRIDSCALE_TF_HEADERS environment variable may be used.
* `request_delay_int` - (Optional) Custom request delay interval in Ms. This time interval is used to delay the synchronous request checks, or delay retryable requests. If omitted, the GRIDSCALE_TF_REQUEST_DELAY_INT environment variable may be used. By default this is set to 1000.
* `max_n_retries` - (Optional) Custom maximum number of retries. The retryable requests can be retried up to max_n_retries. If max_n_retries is reached and the request is not successful, the last error is returned. If omitted, the GRIDSCALE_TF_MAX_N_RETRIES environment variable may be used. By default this is set to 5.