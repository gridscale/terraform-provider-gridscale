package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-gridscale/gridscale"
)

const (
	serverStateFile = "./server.gsstate"
)

func main() {
	serverList := gridscale.NewGlobalServerPowerStateList()
	err := serverList.LoadFromFile(serverStateFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := serverList.FlushToFile(serverStateFile)
		if err != nil {
			panic(err)
		}
	}()
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gridscale.Provider})
}
