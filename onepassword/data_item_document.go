package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceItemDocument() *schema.Resource {
	return &schema.Resource{
		Read: resourceItemDocumentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Document name.",
			},
			"vault": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Vault for document.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "File name.",
			},
			"content": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "File content.",
			},
		},
	}
}
