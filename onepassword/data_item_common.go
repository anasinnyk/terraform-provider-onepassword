package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemCommon() *schema.Resource {
	s := resourceItemCommon().Schema
	s["template"].Required = false
	s["template"].Optional = true

	return &schema.Resource{
		Read:   resourceItemCommonRead,
		Schema: s,
	}
}
