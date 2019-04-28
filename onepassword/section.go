package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"reflect"
	"strings"
)

var sectionSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
		},
		"field": {
			Type:        schema.TypeList,
			ForceNew:    true,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
					},
					"string": {
						Type:        schema.TypeString,
						Optional:    true,
						ForceNew:    true,
						ConflictsWith: []string{
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"url": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: urlValidate,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"phone": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"reference": {
						Type:        schema.TypeString,
						ForceNew:    true,
						Optional:    true,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
						},
					},
					"sex": {
						Type:         schema.TypeString,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: orEmpty(validation.StringInSlice([]string{"female", "male"}, false)),
						ConflictsWith: []string{
							"section.field.string",
							"section.field.phone",
							"section.field.url",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"card_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: orEmpty(validation.StringInSlice([]string{
							"mc",
							"visa",
							"amex",
							"diners",
							"carteblanche",
							"discover",
							"jcb",
							"maestro",
							"visaelectron",
							"laser",
							"unionpay"
							}, false)),
						ConflictsWith: []string{
							"section.field.string",
							"section.field.phone",
							"section.field.url",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.reference",
						},
					},
					"email": {
						Type:         schema.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: emailValidate,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"date": {
						Type:     schema.TypeInt,
						Optional: true,
						ForceNew: true,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"month_year": {
						Type:     schema.TypeInt,
						Optional: true,
						ForceNew: true,
						Description: "Item login section field value for month year.",
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.totp",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"totp": {
						Type:      schema.TypeString,
						Optional:  true,
						Sensitive: true,
						ForceNew:  true,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.concealed",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"concealed": {
						Type:        schema.TypeString,
						Optional:    true,
						Sensitive:   true,
						ForceNew:    true,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.address",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
					},
					"address": {
						Type:        schema.TypeMap,
						Optional:    true,
						ForceNew:    true,
						ConflictsWith: []string{
							"section.field.string",
							"section.field.url",
							"section.field.phone",
							"section.field.email",
							"section.field.date",
							"section.field.month_year",
							"section.field.totp",
							"section.field.concealed",
							"section.field.sex",
							"section.field.card_type",
							"section.field.reference",
						},
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"country": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
								},
								"city": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
								},
								"region": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
								},
								"state": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
								},
								"street": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
								},
								"zip": {
									Type:        schema.TypeString,
									Optional:    true,
									ForceNew:    true,
								},
							},
						},
					},
				},
			},
		},
	},
}

func ParseTags(d *schema.ResourceData) []string {
	tSrc := d.Get("tags").([]interface{})
	tags := make([]string, 0, len(tSrc))
	for _, tag := range tSrc {
		tags = append(tags, tag.(string))
	}
	return tags
}

func ParseFields(s map[string]interface{}) []SectionField {
	fields := []SectionField{}
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
				case "sex":
					f.Type = TypeSex
				case "totp":
					f.Type = TypeConcealed
					f.N = "TOTP_" + f.N
				case "month_year":
					f.Type = TypeMonthYear
				case "url":
					f.Type = TypeURL
				case "card_type":
					f.Type = TypeCard
				default:
					f.Type = SectionFieldType(key)
				}
			}
		}
		fields = append(fields, f)
	}
	return fields
}

func ParseSections(d *schema.ResourceData) []Section {
	sections := []Section{}
	for _, section := range d.Get("section").([]interface{}) {
		s := section.(map[string]interface{})
		sections = append(sections, Section{
			Title:  s["name"].(string),
			Name:   "Section_" + fieldNumber(),
			Fields: ParseFields(s),
		})
	}
	return sections
}
