package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
)

func resourceItemLogin() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemLoginRead,
		Create: resourceItemLoginCreate,
		Delete: resourceItemLoginDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemLoginRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login name.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login username.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "Item login password.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login notes.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Item login tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Vault for item login.",
			},
			"section": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Item login section.",
				Elem:        sectionSchema,
			},
			"url": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "URL for item login.",
				ValidateFunc: urlValidate,
			},
		},
	}
}

func resourceItemLoginRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	log.Printf("[DEBUG] %v", v)
	if err != nil {
		return err
	}

	d.SetId(v.Uuid)
	d.Set("name", v.Overview.Title)
	if err := d.Set("url", v.Overview.Url); err != nil {
		return err
	}
	d.Set("notes", v.Details.Notes)
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)
	for _, field := range v.Details.Fields {
		if field.Name == "username" {
			d.Set("username", field.Value)
		}
		if field.Name == "password" {
			d.Set("password", field.Value)
		}
	}
	sections := make([]map[string]interface{}, 0, len(v.Details.Sections))
	for _, section := range v.Details.Sections {
		fields := make([]map[string]interface{}, 0, len(section.Fields))
		for _, field := range section.Fields {
			f := map[string]interface{}{
				"name": field.Text,
			}
			var key string
			switch field.Type {
			case TypeURL:
				key = "url"
			case TypeMonthYear:
				key = "month_year"
			case TypeConcealed:
				if strings.HasPrefix(field.N, "TOTP_") {
					key = "totp"
				} else {
					key = "concealed"
				}
			default:
				key = string(field.Type)
			}
			f[key] = field.Value
			fields = append(fields, f)
		}
		sections = append(sections, map[string]interface{}{
			"name":  section.Title,
			"field": fields,
		})
	}
	return d.Set("section", sections)
}

func resourceItemLoginCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceItemLoginDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
