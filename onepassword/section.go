package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
)

var sectionSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "Item login section name.",
		},
		"field": {
			Type:        schema.TypeSet,
			Computed:    true,
			Optional:    true,
			Description: "Item login section field.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Item login section name.",
					},
					"string": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Item login section value for string.",
						// ConflictsWith: []string{"url", "phone", "email", "date", "month_year", "totp", "concealed", "address"},
					},
					"url": {
						Type:         schema.TypeString,
						Computed:     true,
						Optional:     true,
						ValidateFunc: urlValidate,
						Description:  "Item login section value for url.",
						// ConflictsWith: []string{"string", "phone", "email", "date", "month_year", "totp", "concealed", "address"},
					},
					"phone": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Description: "Item login section value for phone.",
						// ConflictsWith: []string{"string", "url", "email", "date", "month_year", "totp", "concealed", "address"},
					},
					"email": {
						Type:         schema.TypeString,
						Computed:     true,
						Optional:     true,
						ValidateFunc: emailValidate,
						Description:  "Item login section value for email.",
						// ConflictsWith: []string{"string", "url", "phone", "date", "month_year", "totp", "concealed", "address"},
					},
					"date": {
						Type:     schema.TypeInt,
						Computed: true,
						Optional: true,
						// ValidateFunc: dateValidate,
						Description: "Item login section value for date.",
						// ConflictsWith: []string{"string", "url", "phone", "email", "month_year", "totp", "concealed", "address"},
					},
					"month_year": {
						Type:     schema.TypeInt,
						Computed: true,
						Optional: true,
						// ValidateFunc: monthYearValidate,
						Description: "Item login section value for month year.",
						// ConflictsWith: []string{"string", "url", "phone", "email", "date", "totp", "concealed", "address"},
					},
					"totp": {
						Type:      schema.TypeString,
						Computed:  true,
						Optional:  true,
						Sensitive: true,
						// ValidateFunc: totpValidate,
						Description: "Item login section value for totp.",
						// ConflictsWith: []string{"string", "url", "phone", "email", "date", "month_year", "concealed", "address"},
					},
					"concealed": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						Sensitive:   true,
						Description: "Item login section value for password.",
						// ConflictsWith: []string{"string", "url", "phone", "email", "date", "month_year", "totp", "address"},
					},
					"address": {
						Type:        schema.TypeMap,
						Computed:    true,
						Optional:    true,
						Description: "Item login section value for address.",
						// ConflictsWith: []string{"string", "url", "phone", "email", "date", "month_year", "totp", "concealed"},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"country": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									Description: "Item login section value for address - country.",
								},
								"city": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									Description: "Item login section value for address - city.",
								},
								"region": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									Description: "Item login section value for address - region.",
								},
								"state": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									Description: "Item login section value for address - state.",
								},
								"street": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									Description: "Item login section value for address - street.",
								},
								"zip": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									Description: "Item login section value for address - zip.",
								},
							},
						},
					},
				},
			},
		},
	},
}
