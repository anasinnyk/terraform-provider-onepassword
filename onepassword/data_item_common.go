package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemCommon() *schema.Resource {
	s := resourceItemCommon().Schema
	s["template"].Required = false
	s["template"].Optional = true

	return &schema.Resource{
		ReadContext: resourceItemCommonRead,
		Schema:      s,
	}
}
