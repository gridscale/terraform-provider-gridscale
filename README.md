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

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/bitbucket/gridscale/terraform-provider-gridscale`

```sh
$ mkdir -p $GOPATH/src/bitbucket.org/gridscale; cd $GOPATH/src/bitbucket.org/gridscale
$ git clone git@bitbucket.org:gridscale/terraform-provider-gridscale.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/bitbucket.org/gridscale/terraform-provider-gridscale
$ make build
```

Using the provider
----------------------
See the [`website/docs`](website/docs) directory in this repo to get started on using the gridscale provider, be sure to read [`website/docs/index.html.markdown`](website/docs/index.html.markdown) first. Documentation on how to create resources like servers, storages and networks can be found in [`website/docs/r`](website/docs/r). Documentation on how to add resources like storages, networks and IP addresses to servers, check out the documentation on datasources found in [`website/docs/d`](website/docs/d).

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

Known Issues
---------------------------
The following issues are known to us:

* Changing the name attribute in a template datasource will not trigger storages using this template to be recreated.
* If a storage has snapshots, terraform can not delete it.
* The autorecovery value of a server can't be changed with Terraform.
* The "make website" and "make website test" commands in the makefile don't work for reasons out of our control.
* Adding a storage as boot device an existing server which already has storages linked to it, will result in a system which is unable to boot.