package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemLogin() *schema.Resource {
	return &schema.Resource{
		Read: resourceItemLoginRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login name.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login username.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "Item login password.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Item login tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Vault for item login.",
			},
			"section": {
				Type:        schema.TypeSet,
				Computed:    true,
				Optional:    true,
				Description: "Item login section.",
				Elem:        sectionSchema,
			},
			"url": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "URL for item login.",
				ValidateFunc: urlValidate,
			},
		},
	}
}
