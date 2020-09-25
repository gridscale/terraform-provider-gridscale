# Release Checklist

## Prerequisites

You can only do this with a user that has access to the gridscale organization.

Make sure you have the release signing private key installed (`4841EC2F6BC7BD4515F60C10047EC899C2DC3656`). You also need to export following variables in your environment.

    $ export GITHUB_TOKEN=
    $ export GPG_FINGERPRINT=oss@gridscale.io

Create a new GitHub [personal access token](https://github.com/settings/tokens) if you don't have one yet.

Make sure the working copy is clean and no untracked, modified, or staged files are present. Also make sure changelog is complete and has a release date set.

Add a tag.

    $ git tag v1.7.0

## Build Release Assets and Create a Draft

Make GPG cache the private key passphrase by signing some arbitrary file.

    $ gpg --armor --detach-sign --local-user "oss@gridscale.io" README.md

Type in the passphrase if asked to. This step ensures that GoReleaser can sign the release assets after building and just uses the passphrase cached in the GPG agent. Remove the signature again.

    $ rm -f README.md.asc

Now build the assets and create a draft release.

    $ goreleaser release --rm-dist

This might take a while. When the build is done, push the tag.

    $ git push --tags

## Do the Release and Verify

Finally hop over to [GitHub](https://github.com/gridscale/terraform-provider-gridscale/releases/) and finish up the draft release. Include the changelog entries in the release message.

That's all. Verify the the latest release is present in the Terraform Registry: [registry.terraform.io/providers/gridscale/gridscale/](https://registry.terraform.io/providers/gridscale/gridscale/latest).
