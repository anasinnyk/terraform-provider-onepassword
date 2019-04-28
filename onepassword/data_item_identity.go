package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemIdentity() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemIdentityRead,
		Schema: resourceItemIdentity().Schema,
	}
}
