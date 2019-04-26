package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemLogin() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemLoginRead,
		Schema: resourceItemLogin().Schema,
	}
}
