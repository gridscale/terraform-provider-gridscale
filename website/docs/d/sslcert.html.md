---
layout: "gridscale"
page_title: "gridscale: TLS/SSL certificate"
sidebar_current: "docs-gridscale-datasource-ssl-certificate"
description: |-
  Gets data of a TLS/SSL certificate resource.
---

# gridscale_ssl_certificate

Get data of a TLS/SSL certificate resource. This can be used to link TLS/SSL certificates to a loadbalancer.

## Example Usage

Using the SSL certificate datasource for the creation of a storage:

```terraform
data "gridscale_ssl_certificate" "ssl-certificate-john"{
  resource_id = "xxxx-xxxx-xxxx-xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required) The UUID of the SSL Certificate.

## Attributes Reference

This resource exports the following attributes:

* `id` - The UUID of the SSL Certificate.
* `name` - The human-readable name of the SSH key.
* `common_name` - The common domain name of the SSL certificate.
* `fingerprints` - Defines a list of unique identifiers generated from the MD5, SHA-1, and SHA-256 fingerprints of the certificate.
    * `md5` - MD5 fingerprint of the certificate.
    * `sha256` - SHA256 fingerprint of the certificate.
    * `sha1` - SHA1 fingerprint of the certificate.
* `labels` - The list of labels.
* `status` - status indicates the status of the object.
* `create_time` - The date and time the object was initially created.
* `change_time` - Defines the date and time of the last object change.
* `not_valid_after` - Defines the date after which the certificate is not valid.
