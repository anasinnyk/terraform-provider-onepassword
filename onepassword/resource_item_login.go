package onepassword

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"crypto/rand"
	"github.com/hashicorp/terraform/helper/schema"
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

func fieldNumber() string {
	b := make([]byte, 16)
	rand.Read(b)
	return strings.ToUpper(fmt.Sprintf("%x", b))
}

func resourceItemLoginCreate(d *schema.ResourceData, meta interface{}) error {
	sections := []Section{}
	for _, section := range d.Get("section").([]interface{}) {
		fields := []SectionField{}
		s := section.(map[string]interface{})
		for _, field := range s["field"].([]interface{}) {
			fl := field.(map[string]interface{})
			f := SectionField{
				Text: fl["name"].(string),
			}
			for key, val := range fl {
				if key == "name" {
					continue
				}

				isNotEmptyString := reflect.TypeOf(val).String() == "string" && val != ""
				isNotEmptyInt := reflect.TypeOf(val).String() == "int" && val != 0
				isNotEmptyAddress := strings.HasPrefix(reflect.TypeOf(val).String(), "map") && len(val.(map[string]interface{})) != 0

				if isNotEmptyString || isNotEmptyInt || isNotEmptyAddress {
					f.N = fieldNumber()
					f.Value = val
					switch key {
					case "totp":
						f.Type = TypeConcealed
						f.N = "TOTP_" + f.N
					case "month_year":
						f.Type = TypeMonthYear
					case "url":
						f.Type = TypeURL
					default:
						f.Type = SectionFieldType(key)
					}
				}

			}
			fields = append(fields, f)
		}
		sections = append(sections, Section{
			Title:  s["name"].(string),
			Fields: fields,
		})
	}

	tSrc := d.Get("tags").([]interface{})
	tags := make([]string, 0, len(tSrc))
	for _, tag := range tSrc {
		tags = append(tags, tag.(string))
	}

	item := &Item{
		Uuid:  d.Id(),
		Vault: d.Get("vault").(string),
		Overview: Overview{
			Title: d.Get("name").(string),
			Url:   d.Get("url").(string),
			Tags:  tags,
		},
		Details: Details{
			Notes: d.Get("notes").(string),
			Fields: []Field{
				Field{
					Name:        "username",
					Designation: "username",
					Value:       d.Get("username").(string),
					Type:        FieldText,
				},
				Field{
					Name:        "password",
					Designation: "password",
					Value:       d.Get("password").(string),
					Type:        FieldPassword,
				},
			},
			Sections: sections,
		},
	}
	m := meta.(*Meta)
	err, _ := m.onePassClient.CreateItem(item, "Login")
	if err != nil {
		return err
	}
	return nil
}

func resourceItemLoginDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteItem(getId(d))
	if err == nil {
		d.SetId("")
	}
	return err
}
