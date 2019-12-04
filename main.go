package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-ohishiizumi/ohishiizumi"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ohishiizumi.Provider})
}
