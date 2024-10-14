# Changelog

## 1.27.0 (Oct 14, 2024)

FEATURES:
- Add log-related parameters in k8s resource.

BUG FIXES:
- Fix cluster using deprecated gsk version causes wrong chosen version in k8s resource. [PR #398](https://github.com/gridscale/terraform-provider-gridscale/pull/398).

## 1.26.0 (Jul 11, 2024)

FEATURES:
- Add OICD support to k8s resource. [PR #351](https://github.com/gridscale/terraform-provider-gridscale/pull/351).

## 1.25.0 (Apr 16, 2024)

FEATURES:
- Add parameter `cluster_traffic_encryption` in `gridscale_k8s` resource (only available for k8s ~> 1.29). [PR #320](https://github.com/gridscale/terraform-provider-gridscale/pull/320).

## 1.24.0 (Apr 02, 2024)

FEATURES:
- Add pgaudit_log related fields to `gridscale_postgresql` resource. [PR #299](https://github.com/gridscale/terraform-provider-gridscale/pull/299) 

## 1.23.2 (Feb 27, 2024)

BUG FIXES:
- Fix failed server deletion due to shutting down/turning off a server which is already off. [PR #275](https://github.com/gridscale/terraform-provider-gridscale/pull/275).

## 1.23.1 (Feb 5, 2024)

IMPROVEMENTS:

- Bump `go` from 1.17 to 1.21
- Bump `hashicorp/go-version` from 1.4.0 to 1.6.0
- Bump `hashicorp/terraform-plugin-sdk/v2` from 2.16.0 to 2.32.0
- Bump `aws/aws-sdk-go` from 1.44.114 to 1.50.10

## 1.23.0 (Dec 12, 2023)

FEATURES:
- Add `griscale_mysql8_0` resource [PR #246](https://github.com/gridscale/terraform-provider-gridscale/issues/246)

IMPROVEMENTS:
- Add deprecation notice for `griscale_mysql` (MySQL 5.7)

## 1.22.0 (Sept 26, 2023)

FEATURES:
- Deprecate `gridscale_memcached` resource.

BUG FIXES:
- Fix issue causing "cannot set some dhcp options to null in network resource" [PR #241](https://github.com/gridscale/terraform-provider-gridscale/pull/241). 
- Upgrade some go libraries to fix security issues.

## 1.21.1 (Jul 17, 2023)

BUG FIXES:
- Fix cannot create a k8s cluster with a version not supporting `rocket_storage` [PR #239](https://github.com/gridscale/terraform-provider-gridscale/pull/239).

## 1.21.0 (Jul 11, 2023)

FEATURES:
- Add proxy protocol support to loadbalancer resource [PR #237](https://github.com/gridscale/terraform-provider-gridscale/pull/237)

BUG FIXES:
- Fix `cluster_cidr`'s immutability issue in k8s resource [PR #238](https://github.com/gridscale/terraform-provider-gridscale/pull/238).

IMPROVEMENTS:
- Update k8s resource docs.

## 1.20.0 (Jun 21, 2023)

FEATURES:
- Add rocket storage support to k8s resource [PR #234](https://github.com/gridscale/terraform-provider-gridscale/pull/234)

IMPROVEMENTS:
- Update k8s resource docs.

## 1.19.0 (May 22, 2023)

BUG FIXES:
- Fix server's `auto_recovery` cannot be updated [PR #230](https://github.com/gridscale/terraform-provider-gridscale/pull/230).
- Correct name of `gridscale_object_storage_bucket` in document menu [PR #232](https://github.com/gridscale/terraform-provider-gridscale/pull/232).

FEATURES:
- Allow `user_data_base64` to be set in server resource [PR #231](https://github.com/gridscale/terraform-provider-gridscale/pull/231).
- Allow `comment` and `user_uuid` to be set in object storage access key resource [PR #232](https://github.com/gridscale/terraform-provider-gridscale/pull/232).

IMPROVEMENTS:
- Renew `kubeconfig` in k8s resource when it expires [PR #229](https://github.com/gridscale/terraform-provider-gridscale/pull/229).

## 1.18.1 (Mar 22, 2023)

BUG FIXES:
- Fix missing `darwin/arm64 build` and `windows/arm64 build`.

## 1.18.0 (Mar 15, 2023)

FEATURES:
- Expose k8s private network [PR #219](https://github.com/gridscale/terraform-provider-gridscale/pull/219).
- Allow to set `cluster_cidr` for k8s cluster [PR #221](https://github.com/gridscale/terraform-provider-gridscale/pull/221).

## 1.17.0 (Jan 5, 2023)

FEATURES:
- Allow to set `hardware_profile_config` in server resource [PR #212](https://github.com/gridscale/terraform-provider-gridscale/pull/215).

## 1.16.2 (Nov 7, 2022)

IMPROVEMENTS:
- Add retry when API backend returns error code 409 via [gsclient-go v3.10.1](https://github.com/gridscale/gsclient-go/releases/tag/v3.10.1).
- Improve documentation.

## 1.16.1 (Oct 19, 2022)

BUG FIXES:

* (Hot fix) Fix type assertion bug caused by non-existent k8s_surge_node_count [PR #212](https://github.com/gridscale/terraform-provider-gridscale/pull/212).

## 1.16.0 (Oct 14, 2022)

FEATURES:
- Allow to set `gsk_version` in k8s resource.
- Allow switching between `release` and `gsk_version` in k8s resource.
- Add surge node feature to k8s resource.

IMPROVEMENTS:
- Add retry when API backend returns error code 424, 500 via [gsclient-go v3.10.0](https://github.com/gridscale/gsclient-go/releases/tag/v3.10.0).

## 1.15.0 (Jul 22, 2022)

FEATURES:
- Add private network support for PaaS resources [PR #99](https://github.com/gridscale/terraform-provider-gridscale/pull/199)

IMPROVEMENTS:
- K8S parameter validation is now dynamic [issue #205](https://github.com/gridscale/terraform-provider-gridscale/issues/205)
- Correct names of some resources/datasources in the sidemenu [issue #200](https://github.com/gridscale/terraform-provider-gridscale/issues/200) & [issue #207](https://github.com/gridscale/terraform-provider-gridscale/issues/207).

## 1.14.3 (Apr 7, 2022)

BUG FIXES:

* (Hot fix) Fix updating a network's `name`, `labels` and `l2security` locks a server which is attached to that network until the update is finished. The issue is fixed in [gsclient-go v3.9.1](https://github.com/gridscale/gsclient-go/releases/tag/v3.9.1).

## 1.14.2 (Mar 30, 2022)

IMPROVEMENTS:
- Suppress server-power-state update's error 400. [#196](https://github.com/gridscale/terraform-provider-gridscale/pull/196)

BUG FIXES:
- Fix incorrect maximum number of networks that can be connected to a server. [#197](https://github.com/gridscale/terraform-provider-gridscale/pull/197)

## 1.14.1 (Feb 10, 2022)

IMPROVEMENTS:
- Allow backend to set hardware profile when hardware_profile is not set (q35 at the time of writing) [#191](https://github.com/gridscale/terraform-provider-gridscale/pull/191)
- Allow to set backup location when creating a backup schedule. [#193](https://github.com/gridscale/terraform-provider-gridscale/pull/193)

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
