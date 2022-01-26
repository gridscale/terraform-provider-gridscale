# Changelog

## 1.14.0 (Jan 20, 2022)

FEATURES:
- Add filesystem resource [#189](https://github.com/gridscale/terraform-provider-gridscale/pull/189).
- Allow assigning DHCP IP to a server [#190](https://github.com/gridscale/terraform-provider-gridscale/pull/190).

IMPROVEMENTS:
- Update network's docs.
- Update server's docs.

BUG FIXES:
- Fix network ordering of a server issue [#142](https://github.com/gridscale/terraform-provider-gridscale/issues/142). **NOTE**: The network ordering of the server now corresponds to the order of the networks in the server resource block. The field `ordering` is deprecated. 
- Fix the field `dhcp_range` is not set automatically by the backend, when it is not set manually in terraform [#190](https://github.com/gridscale/terraform-provider-gridscale/pull/190).
- Updating server-storage relations and server-network relations' properties (except the order) will no longer require the server to be shutdown [#190](https://github.com/gridscale/terraform-provider-gridscale/pull/190).

## 1.13.0 (Sept 9, 2021)

FEATURES:
- Add DHCP properties/options to the network resource/data source [#180](https://github.com/gridscale/terraform-provider-gridscale/pull/180).

IMPROVEMENTS:
- Update the network's docs.

## 1.12.0 (Aug 12, 2021)

FEATURES:
- Allow setting storage variant [#171](https://github.com/gridscale/terraform-provider-gridscale/pull/171).

## 1.11.0 (Jun 24, 2021)

IMPROVEMENTS:
- Add more upfront validations [#170](https://github.com/gridscale/terraform-provider-gridscale/pull/170).
- Add field `backup_retention` to MS SQL server backup [#170](https://github.com/gridscale/terraform-provider-gridscale/pull/170).
- Add field `host` to PaaS resource and PaaS-type resources (e.g. postgres, MySQL, etc.) [#169](https://github.com/gridscale/terraform-provider-gridscale/pull/169).

FEATURES:
- Add MySQL resource [#166](https://github.com/gridscale/terraform-provider-gridscale/pull/166).
- Add Memcached resource [#167](https://github.com/gridscale/terraform-provider-gridscale/pull/167).
- Add Redis store/cache resources [#168](https://github.com/gridscale/terraform-provider-gridscale/pull/168).

BUG FIXES:
- Fix MS SQL server backup's params cannot be updated [#170](https://github.com/gridscale/terraform-provider-gridscale/pull/170).

## 1.10.0 (Jun 01, 2021)

IMPROVEMENTS:
- Add more upfront validations [#162](https://github.com/gridscale/terraform-provider-gridscale/pull/162).
- Allow customizing request delay interval and maximum number of retries [#157](https://github.com/gridscale/terraform-provider-gridscale/pull/157).

FEATURES:
- Add TLS/SSL certificate resource and data source [#156](https://github.com/gridscale/terraform-provider-gridscale/pull/156).
- Add MS SQL server resource [#161](https://github.com/gridscale/terraform-provider-gridscale/pull/161).
- Add MariaDB resource [#164](https://github.com/gridscale/terraform-provider-gridscale/pull/164).

BUG FIXES:
- Fix Loadbalancer rs/ds docs (add missing descriptions, fix typo, etc.) [#156](https://github.com/gridscale/terraform-provider-gridscale/pull/156).

## 1.9.1 (Apr 21, 2021)

IMPROVEMENTS:

* Update gsclient-go package to v3.6.2 (which fixes the `PROTOCOL_ERROR` when running terraform in moderate amount of concurrency [#199](https://github.com/gridscale/gsclient-go/pull/199)).

## 1.9.0 (Apr 15, 2021)

IMPROVEMENTS:

* Update gsclient-go package to v3.6.1.

FEATURES:

* Add gridscale k8s resource ([#120](https://github.com/gridscale/terraform-provider-gridscale/issues/120)).
* Add gridscale postgreSQL resource ([#133](https://github.com/gridscale/terraform-provider-gridscale/issues/133)).

## 1.8.4 (Mar 16, 2021)

IMPROVEMENTS:

* Update go version to v1.16.
* Update gsclient-go package to v3.5.0 ([#130](https://github.com/gridscale/terraform-provider-gridscale/pull/131)).
* Use go v1.16 in GH action & enable macOS arm64 build ([#129](https://github.com/gridscale/terraform-provider-gridscale/issues/129)).

## 1.8.3 (Feb 18, 2021)

BUG FIXES:

* Fix inconsistency issue of PaaS service resource, when the `paas_service_template_uuid` is updated outside of terraform ([#123](https://github.com/gridscale/terraform-provider-gridscale/issues/123)).
* Fix `gridscale_backupschedule` and `gridscale_snapshotschedule` resources force to update `next_runtime` when updating other fields ([#126](https://github.com/gridscale/terraform-provider-gridscale/issues/126))

## 1.8.2 (Jan 26, 2021)

BUG FIXES:

* Reading SSH public key from file via `file()` function will not cause key to be set in every apply ([#116](https://github.com/gridscale/terraform-provider-gridscale/issues/116)).

## 1.8.1 (Jan 21, 2021)

IMPROVEMENTS:

* We fixed the User-Agent header that is sent with Terraform requests to something useful ([#108](https://github.com/gridscale/terraform-provider-gridscale/issues/108)).

BUG FIXES:

* PaaS services are not destroyed and re-created anymore when a service template changes ([#109](https://github.com/gridscale/terraform-provider-gridscale/issues/109)).
* SSH Public key can be added via function `file()` without any problems with whitespace ([#112](https://github.com/gridscale/terraform-provider-gridscale/issues/112)).

## 1.8.0 (Jan 05, 2021)

IMPROVEMENTS:

* Update docs.
* Update gsclient-go package to v3.3.1.

BUG FIXES:

* Fix storage_type is not set when cloning a storage See [#105](https://github.com/gridscale/terraform-provider-gridscale/issues/105)

FEATURES:

* Add storage import (from storage backups) feature.

## 1.7.4 (Nov 03, 2020)

IMPROVEMENTS:

* Reword docs.
* Add an example (and explanation) about firewall rules in server-network relation.
* Explain how ordering of network interfaces works.

BUG FIXES:

* Fix ordering of network interfaces on the host is NOT the same as defined in the Terraform definition (top-down order). See [#99](https://github.com/gridscale/terraform-provider-gridscale/issues/99).
* Enable firewall only when at least one firewall rule is set. In previous version, when no firewall rules are set, the default firewall rules are added. This makes all ports blocked. See [#100](https://github.com/gridscale/terraform-provider-gridscale/issues/100)

## 1.7.3 (Nov 02, 2020)

BUG FIXES:

* The ordering of networks in a server's relation now can be set. See [#95](https://github.com/gridscale/terraform-provider-gridscale/issues/95).

## 1.7.2 (Oct 29, 2020)

CHANGES:

* Update gsclient-go package to v3.2.2.
* Allow to omit user UUID and API token in requests when they are empty.
* Update release checklist. No need to do GPG signing and building manually. All done by the pipeline now.

## 1.7.1 (Oct 15, 2020)

BUG FIXES:

* The provider is now applying default inbound firewall rules. See [#89](https://github.com/gridscale/terraform-provider-gridscale/issues/89).
* Fix turning off a server even when it is already shutdown.

## 1.7.0 (Sept 11, 2020)

FEATURES:

* Support marketplace application features.
* Support storage backup functionality and schedule storage backup.

IMPROVEMENTS:

* Update gsclient-go package to v3.2.1.
* Replace Travis CI with GitHub Actions.

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
