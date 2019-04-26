package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGroupRead,
		Schema: resourceGroup().Schema,
	}
}
