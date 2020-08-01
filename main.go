package main

import (
	"github.com/anasinnyk/terraform-provider-1password/onepassword"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return onepassword.Provider()
		},
	})
}
