package onepassword

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func dataSourceItemDocument() *schema.Resource {
	s := resourceItemDocument().Schema
	delete(s, "file_path")

	return &schema.Resource{
		ReadContext: resourceItemDocumentRead,
		Schema:      s,
	}
}
