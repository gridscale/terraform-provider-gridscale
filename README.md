# Terraform gridscale Provider

[![Build status](https://github.com/gridscale/terraform-provider-gridscale/workflows/Test/badge.svg)](https://github.com/gridscale/terraform-provider-gridscale/actions)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/gridscale/terraform-provider-gridscale?label=release)](https://github.com/gridscale/terraform-provider-gridscale/releases)
[![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px">

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

| Feature | Availability |
|---|:---:|
| Server (CRUD) | :heavy_check_mark: |
| Server dependency (link/unlink) | :heavy_check_mark: |
| Load balancer (CRUD) | :heavy_check_mark: |
| PaaS (CRUD) | :heavy_check_mark: |
| Storage (CRUD) | :heavy_check_mark: |
| Object storage (CRUD) | :heavy_check_mark: |
| IP address (CRUD) | :heavy_check_mark: |
| Network (CRUD) | :heavy_check_mark: |
| Security zone (PaaS) (CRUD) | :heavy_check_mark: |
| Firewall (CRUD) | :heavy_check_mark: |
| SSH key (CRUD) | :heavy_check_mark: |
| ISO Image (CRUD) | :heavy_check_mark: |
| Snapshot (CRUD) | :heavy_check_mark: |
| Snapshot rollback | :heavy_check_mark: |
| Snapshot to S3 | :heavy_check_mark: |
| Snapshot schedule (CRUD) | :heavy_check_mark: |
| Template (CRUD) | :heavy_check_mark: |
| Multiple project support ([Workaround](https://github.com/gridscale/terraform-examples/tree/master/multi-project)) | :x: |
| Marketplace | :heavy_check_mark: |

## Development

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

    $ make build

To test the provider, simply run

    $ make test

In order to run the full suite of acceptance tests, run `make testacc`. You will need to set `GRIDSCALE_UUID`, `GRIDSCALE_TOKEN`, and `GRIDSCALE_URL` environment variables to point to an existing project when running acceptance tests.

*Note:* acceptance tests create real resources, and often cost money to run.

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
- If a storage has snapshots, terraform can not delete it.
- The autorecovery value of a server can't be changed with Terraform.
