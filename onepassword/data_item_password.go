package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemPassword() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemPasswordRead,
		Schema: resourceItemPassword().Schema,
	}
}
