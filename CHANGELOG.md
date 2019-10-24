## 1.3.0 (October 25, 2019)

FEATURES:

* Add datasource for all available resource

IMPROVEMENTS:

* Switch to gsclient-go v2.0.0 from github
* Fix issue: Terraform destroy raises error when instances powered up (https://github.com/terraform-providers/terraform-provider-gridscale/issues/13)
* Fix issue: Reducing cores / memory does not cause server shutdown (https://github.com/terraform-providers/terraform-provider-gridscale/issues/12)
* Add tests for all available datasource
* Fix all datasources missing `Schema`

## 1.2.0 (July 30, 2019)

FEATURES:

* Add support for LBaaS (CH-15)


## 1.1.0 (July 10, 2019)

FEATURES:

* Assure compatibility with terraform 0.12
* Allow using availability zone C (GH-10)

IMPROVEMENTS:

* Switch to gsclient-go from github (GH-14)

## 1.0.0 (April 30, 2019)

* Initial release.
