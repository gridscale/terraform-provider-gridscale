## 1.7.2 (Oct 28, 2020)

IMPROVEMENTS:

* Update gsclient-go package to v3.2.2.
* Allow to omit uuid and token when they are empty.
* Update release checklist (No need to do pgp signing and build binary files).

## 1.7.1 (Oct 15, 2020)

BUG FIXES:

* The provider is now applying default inbound firewall rules. See [#89](https://github.com/gridscale/terraform-provider-gridscale/issues/89).
* Fix turning off a server even when it is already shutdown.

## 1.7.0 (Sept 11, 2020)

FEATURES:

* Support marketplace application features.
* Support storage backup functionalities and schedule storage backup.

IMPROVEMENTS:

* Update gsclient-go package to v3.2.1.
* Replace travis with github actions.

BUG FIXES:

* Fix bug causing `next_runtime` fields of snapshot schedule and backup schedule to be changed by gs server unexpectedly.

## 1.6.3 (Aug 18, 2020)

Prepare publishing to Terraform Registry.

IMPROVEMENTS:

* Remove redundant types in data sources.
* Size and type of a storage can be modified.

BUG FIXES:

* Update of storage type won't force to create new storage.

## 1.6.2 (July 07, 2020)

IMPROVEMENTS:

* Custom HTTP headers are supported.

## 1.6.1 (June 30, 2020)

IMPROVEMENTS:

* Update gsclient-go package to v3.1.0
* Update and tidy the vendor directory.

## 1.6.0 (June 02, 2020)

FEATURES:

* Support exporting snapshot to object storage.
* Support specific timeouts.

IMPROVEMENTS:

* Update gsclient-go package to v3.0.0
* Update and tidy the vendor directory.
* Remove unnecessary/dummy variables.
* Skip 404 when deleting a resource (and 409 when deleting a server-related resource).
* Reconstruct some internal packages (rename/create).
* Increase default timeouts of PaaS operations to 15 minutes.

## 1.5.1 (April 06, 2020)

IMPROVEMENTS:

* Fix documentation (wrong navigation, reorder the side menu).
* Add link to multi-project example.
* Add configurable timeout.

## 1.5.0 (January 24, 2020)

FEATURES:

* Support all available gridscale resources
* Improve documentation
* Server CRUD is faster

IMPROVEMENTS:

* Switch to gsclient-go v2.2.1 from GitHub (better connection error handling)
* Handle all errors when setting values
* Robust error reporting
* Fix bugs caused by:
  * Missing properties
  * Type mismatch
  * Weak error handling

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

* Switch to gsclient-go v2.0.0 from GitHub
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

* Switch to gsclient-go from GitHub ([#14](https://github.com/terraform-providers/terraform-provider-template/issues/14))

## 1.0.0 (April 30, 2019)

* Initial release.
