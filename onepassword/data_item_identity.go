package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemIdentity() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceItemIdentityRead,
		Schema:      resourceItemIdentity().Schema,
	}
}
