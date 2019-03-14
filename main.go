package main

import (
	"github.com/gridscale/terraform-provider-gridscale/gridscale"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gridscale.Provider})
}
