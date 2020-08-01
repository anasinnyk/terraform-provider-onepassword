package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemSecureNote() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceItemSecureNoteRead,
		Schema:      resourceItemSecureNote().Schema,
	}
}
