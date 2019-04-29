package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemSoftwareLicense() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemSoftwareLicenseRead,
		Schema: resourceItemSoftwareLicense().Schema,
	}
}
