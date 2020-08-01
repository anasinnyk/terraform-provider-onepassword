package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemCreditCard() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceItemCreditCardRead,
		Schema:      resourceItemCreditCard().Schema,
	}
}
