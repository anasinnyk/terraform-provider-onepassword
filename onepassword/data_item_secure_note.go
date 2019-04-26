package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemSecureNote() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemSecureNoteRead,
		Schema: resourceItemSecureNote().Schema,
	}
}
