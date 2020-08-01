package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemLogin() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceItemLoginRead,
		Schema:      resourceItemLogin().Schema,
	}
}
