package main

import (
	"github.com/hashicorp/terraform/plugin"
	"./gridscale"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gridscale.Provider})
}
