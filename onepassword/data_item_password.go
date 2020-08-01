package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemPassword() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceItemPasswordRead,
		Schema:      resourceItemPassword().Schema,
	}
}
