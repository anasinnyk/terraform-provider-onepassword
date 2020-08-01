package onepassword

import (
	"reflect"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func sectionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"field": {
				Type:     schema.TypeList,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"string": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"url": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: urlValidate,
						},
						"phone": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"reference": {
							Type:     schema.TypeString,
							ForceNew: true,
							Optional: true,
						},
						"sex": {
							Type:         schema.TypeString,
							ForceNew:     true,
							Optional:     true,
							ValidateFunc: orEmpty(validation.StringInSlice([]string{"female", "male"}, false)),
						},
						"card_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
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
								"unionpay",
							}, false)),
						},
						"email": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: emailValidate,
						},
						"date": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"month_year": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "Item login section field value for month year.",
						},
						"totp": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"concealed": {
							Type:      schema.TypeString,
							Optional:  true,
							Sensitive: true,
							ForceNew:  true,
						},
						"address": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
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
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"city": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"region": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"state": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"street": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"zip": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ParseTags(d *schema.ResourceData) []string {
	tSrc := d.Get("tags").([]interface{})
	tags := make([]string, 0, len(tSrc))
	for _, tag := range tSrc {
		tags = append(tags, tag.(string))
	}
	return tags
}

func ParseField(fl map[string]interface{}) SectionField {
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
			f.N = f.Text
			if val, err := fieldNumber(); err == nil {
				f.N = val
			}
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
	return f
}

func ParseFields(s map[string]interface{}) []SectionField {
	fields := []SectionField{}
	for _, field := range s["field"].([]interface{}) {
		fl := field.(map[string]interface{})
		fields = append(fields, ParseField(fl))
	}
	return fields
}

func ParseSections(d *schema.ResourceData) []Section {
	sections := []Section{}
	for _, section := range d.Get("section").([]interface{}) {
		s := section.(map[string]interface{})
		secName := s["name"].(string)
		if val, err := fieldNumber(); err == nil {
			secName = val
		}
		sections = append(sections, Section{
			Title:  s["name"].(string),
			Name:   "Section_" + secName,
			Fields: ParseFields(s),
		})
	}
	return sections
}
