package main

import (
	"bitbucket.org/gridscale/terraform-provider-gridscale/gridscale"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gridscale.Provider})
}
