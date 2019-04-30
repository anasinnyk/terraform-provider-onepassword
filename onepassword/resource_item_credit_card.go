package onepassword

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceItemCreditCard() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemCreditCardRead,
		Create: resourceItemCreditCardCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemCreditCardRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"main": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cardholder": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"type": {
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
						"number": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cvv": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"expiry_date": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"valid_from": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"field": sectionSchema().Schema["field"],
					},
				},
			},
			"section": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     sectionSchema(),
			},
		},
	}
}

func resourceItemCreditCardRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(CreditCardCategory) {
		return errors.New("item is not from " + string(CreditCardCategory))
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Overview.Title); err != nil {
		return err
	}
	if err := d.Set("tags", v.Overview.Tags); err != nil {
		return err
	}
	if err := d.Set("vault", v.Vault); err != nil {
		return err
	}
	if err := d.Set("notes", v.Details.Notes); err != nil {
		return err
	}
	return parseSectionFromSchema(v.Details.Sections, d, []SectionGroup{
		SectionGroup{
			Name:     "main",
			Selector: "",
			Fields: map[string]string{
				"cardholder":  "cardholder",
				"number":      "ccnum",
				"type":        "type",
				"cvv":         "cvv",
				"expiry_date": "expiry",
				"valid_from":  "validFrom",
			},
		},
	})
}

func resourceItemCreditCardCreate(d *schema.ResourceData, meta interface{}) error {
	main := d.Get("main").([]interface{})[0].(map[string]interface{})
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(CreditCardCategory),
		Details: Details{
			Notes: d.Get("notes").(string),
			Sections: append(
				[]Section{
					Section{
						Title: main["title"].(string),
						Name:  "",
						Fields: append([]SectionField{
							SectionField{
								Type:  "string",
								Text:  "cardholder",
								Value: main["cardholder"].(string),
								N:     "cardholder",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
									"keyboard":           "Default",
								},
							},
							SectionField{
								Type:  "cctype",
								Text:  "type",
								Value: main["type"].(string),
								N:     "type",
								A: Annotation{
									guarded: "yes",
								},
							},
							SectionField{
								Type:  "string",
								Text:  "number",
								Value: main["number"].(string),
								N:     "ccnum",
								A: Annotation{
									guarded:         "yes",
									clipboardFilter: "0123456789",
								},
								Inputs: map[string]string{
									"keyboard": "NumberPad",
								},
							},
							SectionField{
								Type:  "concealed",
								Text:  "verification number",
								Value: main["cvv"].(string),
								N:     "cvv",
								A: Annotation{
									guarded:  "yes",
									generate: "off",
								},
								Inputs: map[string]string{
									"keyboard": "NumberPad",
								},
							},
							SectionField{
								Type:  "monthYear",
								Text:  "expiry date",
								Value: main["expiry_date"].(int),
								N:     "expiry",
								A: Annotation{
									guarded: "yes",
								},
							},
							SectionField{
								Type:  "monthYear",
								Text:  "valid from",
								Value: main["valid_from"].(int),
								N:     "validFrom",
								A: Annotation{
									guarded: "yes",
								},
							},
						}, ParseFields(main)...),
					},
				},
				ParseSections(d)...,
			),
		},
		Overview: Overview{
			Title: d.Get("name").(string),
			Tags:  ParseTags(d),
		},
	}
	m := meta.(*Meta)
	err := m.onePassClient.CreateItem(item)
	if err != nil {
		return err
	}
	d.SetId(item.UUID)
	return nil
}
