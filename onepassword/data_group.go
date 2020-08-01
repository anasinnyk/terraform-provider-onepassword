package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceGroupRead,
		Schema:      resourceGroup().Schema,
	}
}
