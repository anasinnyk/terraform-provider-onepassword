package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/anasinnyk/terraform-provider-1password/onepassword"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: onepassword.Provider,
	})
}
