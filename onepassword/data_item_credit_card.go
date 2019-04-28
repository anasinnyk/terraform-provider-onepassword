package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemCreditCard() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemCreditCardRead,
		Schema: resourceItemCreditCard().Schema,
	}
}
