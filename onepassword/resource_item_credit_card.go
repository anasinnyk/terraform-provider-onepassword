package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceItemCreditCard() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceItemCreditCardRead,
		CreateContext: resourceItemCreditCardCreate,
		DeleteContext: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceItemCreditCardRead(ctx, d, meta); err.HasError() {
					return []*schema.ResourceData{d}, errors.New(err[0].Summary)
				}
				return []*schema.ResourceData{d}, nil
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

func resourceItemCreditCardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return diag.FromErr(err)
	}
	if v.Template != Category2Template(CreditCardCategory) {
		return diag.FromErr(errors.New("item is not from " + string(CreditCardCategory)))
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Overview.Title); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", v.Overview.Tags); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vault", v.Vault); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notes", v.Details.Notes); err != nil {
		return diag.FromErr(err)
	}
	if err := parseSectionFromSchema(v.Details.Sections, d, []SectionGroup{
		{
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
	}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceItemCreditCardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	main := d.Get("main").([]interface{})[0].(map[string]interface{})
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(CreditCardCategory),
		Details: Details{
			Notes: d.Get("notes").(string),
			Sections: append(
				[]Section{
					{
						Title: main["title"].(string),
						Name:  "",
						Fields: append([]SectionField{
							{
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
							{
								Type:  "cctype",
								Text:  "type",
								Value: main["type"].(string),
								N:     "type",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
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
							{
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
							{
								Type:  "monthYear",
								Text:  "expiry date",
								Value: main["expiry_date"].(int),
								N:     "expiry",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
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
		return diag.FromErr(err)
	}
	d.SetId(item.UUID)
	return nil
}
