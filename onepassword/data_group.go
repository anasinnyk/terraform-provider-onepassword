package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read: resourceGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Optional:    true,
				Description: "Group name.",
			},
		},
	}
}
