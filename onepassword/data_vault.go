package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceVault() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceVaultRead,
		Schema:      resourceVault().Schema,
	}
}
