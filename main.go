package main

import (
	"github.com/hashicorp/terraform/plugin"
	"wouter_gridscale@bitbucket.org/wouter_gridscale/terraform-provider-gridscale/gridscale"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: template.Provider})
}
