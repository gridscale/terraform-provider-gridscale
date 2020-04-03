Terraform gridscale Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Terraform team at [gridscale](https://www.gridscale.io/).

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- [Go](https://golang.org/doc/install) 1.13.x (to build the provider plugin)

Building the Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/gridscale/terraform-provider-gridscale`

```sh
$ mkdir -p $GOPATH/src/github.com/gridscale; cd $GOPATH/src/github.com/gridscale
$ git clone git@github.com:gridscale/terraform-provider-gridscale.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/gridscale/terraform-provider-gridscale
$ make build
```

Using the Provider
------------------

See the [gridscale provider documentation](https://www.terraform.io/docs/providers/gridscale) to get started on using the gridscale provider.

Alternatively, this documentation can also be found within this repository. Check out [`website/docs/index.html.markdown`](website/docs/index.html.markdown) to get started. Documentation on how to create resources like servers, storages and networks can be found in [`website/docs/r`](website/docs/r). Documentation on how to add resources like storages, networks and IP addresses to servers, check out the documentation on datasources found in [`website/docs/d`](website/docs/d).

Available Features
---------------------------

| Feature | Availability |
|---|:---:|
| Server (CRUD) | :heavy_check_mark: |  
| Server dependency (Link/Unlink) | :heavy_check_mark: |  
| Loadbalancer (CRUD) | :heavy_check_mark: |  
| PaaS (CRUD) | :heavy_check_mark: |  
| Storage (CRUD) | :heavy_check_mark: |  
| Object Storage (CRUD) | :heavy_check_mark: |  
| IP address (CRUD) | :heavy_check_mark: |  
| Network (CRUD) | :heavy_check_mark: |  
| Security zone (PaaS) (CRUD) | :heavy_check_mark: |  
| Firewall (CRUD) | :heavy_check_mark: |  
| SSH key (CRUD) | :heavy_check_mark: |  
| ISO Image (CRUD) | :heavy_check_mark: |  
| Snapshot (CRUD) | :heavy_check_mark: |  
| Snapshot rollback | :heavy_check_mark: |  
| Snapshot to S3 | :x: |  
| Snapshot schedule (CRUD) | :heavy_check_mark: |  
| Template (CRUD) | :heavy_check_mark: |  
| Multiple project support ([Workaround](https://github.com/gridscale/terraform-examples/tree/master/multi-project)) | :x: |

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-gridscale
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to run a specific Acceptance test.

```sh
$ make testacc TEST=./gridscale/ TESTARGS='-run=TestAccResourceGridscaleLoadBalancerBasic'
```

Known Issues
---------------------------

The following issues are known to us:

- Changing the name attribute in a template datasource will not trigger storages using this template to be recreated.
- If a storage has snapshots, terraform can not delete it.
- The autorecovery value of a server can't be changed with Terraform.
