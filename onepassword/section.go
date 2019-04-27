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
			Computed:    true,
			ForceNew:    true,
			Optional:    true,
			Description: "Item login section name.",
		},
		"field": {
			Type:        schema.TypeList,
			Computed:    true,
			ForceNew:    true,
			Optional:    true,
			Description: "Item login section field.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Computed:    true,
						ForceNew:    true,
						Optional:    true,
						Description: "Item login section field name.",
					},
					"string": {
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
						ForceNew:    true,
						Description: "Item login section field value for string.",
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
						Computed:     true,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: urlValidate,
						Description:  "Item login section field value for url.",
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
						Computed:    true,
						ForceNew:    true,
						Optional:    true,
						Description: "Item login section field value for phone.",
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
						Computed:    true,
						ForceNew:    true,
						Optional:    true,
						Description: "Item login section field value for reference.",
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
						Computed:     true,
						ForceNew:     true,
						Optional:     true,
						ValidateFunc: orEmpty(validation.StringInSlice([]string{"female", "male"}, false)),
						Description:  "Item login section field value for sex.",
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
						Computed:     true,
						Optional:     true,
						ForceNew:     true,
						Description:  "Item login section field value for card type.",
						ValidateFunc: orEmpty(validation.StringInSlice([]string{"mc", "visa", "amex", "diners", "carteblanche", "discover", "jcb", "maestro", "visaelectron", "laser", "unionpay"}, false)),
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
						Computed:     true,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: emailValidate,
						Description:  "Item login section field value for email.",
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
						Computed: true,
						Optional: true,
						ForceNew: true,
						// ValidateFunc: dateValidate,
						Description: "Item login section field value for date.",
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
						Computed: true,
						Optional: true,
						ForceNew: true,
						// ValidateFunc: monthYearValidate,
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
						Computed:  true,
						Optional:  true,
						Sensitive: true,
						ForceNew:  true,
						// ValidateFunc: totpValidate,
						Description: "Item login section field value for totp.",
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
						Computed:    true,
						Optional:    true,
						Sensitive:   true,
						ForceNew:    true,
						Description: "Item login section field value for password.",
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
						Computed:    true,
						Optional:    true,
						ForceNew:    true,
						Description: "Item login section field value for address.",
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
									Computed:    true,
									Optional:    true,
									ForceNew:    true,
									Description: "Item login section field value for address - country.",
								},
								"city": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									ForceNew:    true,
									Description: "Item login section field value for address - city.",
								},
								"region": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									ForceNew:    true,
									Description: "Item login section field value for address - region.",
								},
								"state": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									ForceNew:    true,
									Description: "Item login section field value for address - state.",
								},
								"street": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									ForceNew:    true,
									Description: "Item login section field value for address - street.",
								},
								"zip": {
									Type:        schema.TypeString,
									Computed:    true,
									Optional:    true,
									ForceNew:    true,
									Description: "Item login section field value for address - zip.",
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

func ParseSections(d *schema.ResourceData) []Section {
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
		sections = append(sections, Section{
			Title:  s["name"].(string),
			Name:   "Section_" + fieldNumber(),
			Fields: fields,
		})
	}
	return sections
}
