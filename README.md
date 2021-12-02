# Terraform gridscale Provider

[Build status](https://github.com/gridscale/terraform-provider-gridscale#available-features)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/gridscale/terraform-provider-gridscale?label=release)](https://github.com/gridscale/terraform-provider-gridscale/releases)
[![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)

## Maintainers

This provider plugin is maintained by the Terraform team at [gridscale](https://www.gridscale.io/).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) ≥ 0.12.x
- [Go](https://golang.org/doc/install) ≥ 1.13 (to build the provider plugin)

## Building the Provider

Clone repository to: `$GOPATH/src/github.com/gridscale/terraform-provider-gridscale`

    $ mkdir -p $GOPATH/src/github.com/gridscale; cd $GOPATH/src/github.com/gridscale
    $ git clone git@github.com:gridscale/terraform-provider-gridscale.git

Enter the provider directory and build the provider

    $ cd $GOPATH/src/github.com/gridscale/terraform-provider-gridscale
    $ make build

## Using the Provider

See the [gridscale provider documentation](https://registry.terraform.io/providers/gridscale/gridscale/latest/docs) to get started on using the gridscale provider.

Alternatively, this documentation can also be found within this repository. Check out [`website/docs/index.html.markdown`](website/docs/index.html.markdown) to get started. Documentation on how to create resources like servers, storages and networks can be found in [`website/docs/r`](website/docs/r). Documentation on how to add resources like storages, networks and IP addresses to servers, check out the documentation on datasources found in [`website/docs/d`](website/docs/d).

## Available Features

| Feature | Availability | Test |
|---|:---:|:---:|
| Server (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/server.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/server.yml) |
| Server dependency (link/unlink) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/server.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/server.yml) |
| Load balancer (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/loadbalancer.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/loadbalancer.yml) |
| PaaS (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/paas.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/paas.yml) |
| K8S (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/k8s.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/k8s.yml) |
| MySQL (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/mysql.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/mysql.yml) |
| MSSQL (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/mssql.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/mssql.yml) |
| MariaDB (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/mariadb.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/mariadb.yml) |
| Postgres (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/postgres.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/postgres.yml) |
| Memcached (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/memcached.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/memcached.yml) |
| Redis cache (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/redis.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/redis.yml) |
| Redis store (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/redis.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/redis.yml) |
| Storage (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/storage.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/storage.yml) |
| SSL Cert (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/sslcert.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/sslcert.yml) |
| Object storage (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/object_storage.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/object_storage.yml) |
| IP address (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/ipv4_ipv6.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/ipv4_ipv6.yml) |
| Network (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/network.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/network.yml) |
| Security zone (PaaS) (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/security_zone.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/security_zone.yml) |
| Firewall (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/firewall.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/firewall.yml) |
| SSH key (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/sshkey.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/sshkey.yml) |
| ISO Image (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/isoimage.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/isoimage.yml) |
| Snapshot (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml) |
| Snapshot rollback | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml) |
| Snapshot to S3 | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml) |
| Snapshot schedule (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/snapshot.yml) |
| Template (CRUD) | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/template.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/template.yml) |
| Multiple project support ([Workaround](https://github.com/gridscale/terraform-examples/tree/master/multi-project)) | :x: |
| Marketplace | :heavy_check_mark: | [![Build status](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/marketplace_app.yml/badge.svg?branch=master)](https://github.com/gridscale/terraform-provider-gridscale/actions/workflows/marketplace_app.yml) |

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

    $ make build

To test the provider, simply run

    $ make test

To run the full suite of acceptance tests execute `make testacc`. You will need to set `GRIDSCALE_UUID`, `GRIDSCALE_TOKEN`, and `GRIDSCALE_URL` environment variables to point to an existing project when running acceptance tests.

*Note:* acceptance tests create real resources and often cost money to run.

    $ make testacc

To run a specific acceptance test, use `TESTARGS`.

    $ make testacc \
        TEST=./gridscale \
        TESTARGS='-run=TestAccResourceGridscaleLoadBalancerBasic'

## Releasing the Provider

A [GoReleaser](https://goreleaser.com/) configuration is provided that produces build artifacts matching the [layout required](https://www.terraform.io/docs/registry/providers/publishing.html#manually-preparing-a-release) to publish the provider in the Terraform Registry.

Releases will appear as drafts. Once marked as published on the GitHub Releases page, they will become available via the Terraform Registry. Releases are signed with key ID (long) `4841EC2F6BC7BD4515F60C10047EC899C2DC3656`.

Jump to the [Release Checklist](release-checklist.md) for details.

## Known Issues

The following issues are known to us:

- Changing the name attribute in a template datasource will not trigger storages using this template to be recreated.
- If a storage has snapshots, Terraform cannot delete it.
- The autorecovery value of a server can't be changed with Terraform.
