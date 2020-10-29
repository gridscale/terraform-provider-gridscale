# Release Checklist

## Add and push a new tag

That is really all you need to do now. Release assets will be build automatically by GoReleaser and the checksums file will be signed with the release signing key. A draft release will be created for the tag.

    $ git tag v1.7.0
    $ git push --tags

## Do the Release and Verify

After the build has finished, hop over to [GitHub](https://github.com/gridscale/terraform-provider-gridscale/releases/) and finish up the draft release. Include the changelog entries in the release message. If you don't find a draft release, you might want to check the GH Action output for a link.

That's all. Verify the the latest release is present in the Terraform Registry: [registry.terraform.io/providers/gridscale/gridscale/](https://registry.terraform.io/providers/gridscale/gridscale/latest).
