## 1.4.1 (Unreleased)
## 1.4.0 (October 31, 2019)

FEATURES:

* Support firewall configuration

IMPROVEMENTS:

* Turn off server synchronously when removing resource attached to it 
* Bootdevice attribute has become `computed`
* firewall_template_uuid has become `optional`
* Server dependency manager features: Create/Update/Remove server's relations.

## 1.3.0 (Unreleased)

FEATURES:

* Add datasource for all available resources

IMPROVEMENTS:

* Switch to gsclient-go v2.0.0 from github
* Fix issue #13: Terraform destroy raises error when instances powered up (https://github.com/terraform-providers/terraform-provider-gridscale/issues/13)
* Fix issue #12: Reducing cores / memory does not cause server shutdown (https://github.com/terraform-providers/terraform-provider-gridscale/issues/12)
* Add tests for all available datasource
* Fix all datasources missing `Schema`
* Update website/docs

## 1.2.0 (July 30, 2019)

FEATURES:

* Add support for LBaaS (CH-15)


## 1.1.0 (July 10, 2019)

FEATURES:

* Assure compatibility with terraform 0.12
* Allow using availability zone C ([#10](https://github.com/terraform-providers/terraform-provider-template/issues/10))

IMPROVEMENTS:

* Switch to gsclient-go from github ([#14](https://github.com/terraform-providers/terraform-provider-template/issues/14))

## 1.0.0 (April 30, 2019)

* Initial release.
