package main

import (
	"github.com/anasinnyk/terraform-provider-1password/onepassword"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: onepassword.Provider,
	})
}
