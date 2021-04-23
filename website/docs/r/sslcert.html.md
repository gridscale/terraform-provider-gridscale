---
layout: "gridscale"
page_title: "gridscale: gridscale_ssl_certificate"
sidebar_current: "docs-gridscale-resource-ssl_certificate"
description: |-
  Manages a SSL Certificate resource in gridscale.
---

# gridscale_ssl_certificate

Provides a SSL Certificate resource. This can be used to create and delete SSL Certificates.
A SSL Certificate can be attached to a loadbalancer.

## Example Usage

The following example shows how one might use this resource to add an SSL Certificate to gridscale:

```terraform
resource "gridscale_ssl_certificate" "ssl-certificate-john"{
  name = "john's computer"
  private_key = "a private-key of the SSL certificate" //or file("/path/to/private.key")
  leaf_certificate = "a public SSL of the SSL certificate" //or file("/path/to/certificate.cert")
  timeouts {
    create="10m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, Force New) The human-readable name of the object. It supports the full UTF-8 character set, with a maximum of 64 characters.

* `private_key` - (Required, Force New) The PEM-formatted private-key of the SSL certificate.

* `leaf_certificate` - (Required, Force New) The PEM-formatted public SSL of the SSL certificate.

* `certificate_chain` - (Optional, Force New) The PEM-formatted full-chain between the certificate authority and the domain's SSL certificate.

* `labels` - (Optional, Force New) List of labels in the format [ "label1", "label2" ].

## Timeouts

Timeouts configuration options (in seconds):
More info: [terraform.io/docs/configuration/resources.html#operation-timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts)

* `create` - (Default value is "5m" - 5 minutes) Used for creating a resource.
* `update` - (Default value is "5m" - 5 minutes) Used for updating a resource.
* `delete` - (Default value is "5m" - 5 minutes) Used for deleting a resource.

## Attributes

This resource exports the following attributes:

* `name` - See Argument Reference above.
* `common_name` - The common domain name of the SSL certificate.
* `ssl_certificate` - See Argument Reference above.
* `leaf_certificate` - See Argument Reference above.
* `certificate_chain` - See Argument Reference above.
* `fingerprints` - Defines a list of unique identifiers generated from the MD5, SHA-1, and SHA-256 fingerprints of the certificate.
    * `md5` - MD5 fingerprint of the certificate.
    * `sha256` - SHA256 fingerprint of the certificate.
    * `sha1` - SHA1 fingerprint of the certificate.
* `labels` - See Argument Reference above.
* `status` - status indicates the status of the object.
* `create_time` - The date and time the object was initially created.
* `change_time` - Defines the date and time of the last object change.
* `not_valid_after` - Defines the date after which the certificate is not valid.
