package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceVault() *schema.Resource {
	return &schema.Resource{
		Read: resourceVaultRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Optional:    true,
				Description: "Vault name.",
			},
		},
	}
}
