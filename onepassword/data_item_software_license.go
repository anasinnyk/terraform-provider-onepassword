package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemSoftwareLicense() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceItemSoftwareLicenseRead,
		Schema:      resourceItemSoftwareLicense().Schema,
	}
}
