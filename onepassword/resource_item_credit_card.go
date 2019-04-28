package onepassword

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
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
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Item credit card name.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Item credit card tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Vault for item credit card.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Item credit card.",
			},
			"main": {
				Type:        schema.TypeList,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Item credit card - main.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:        schema.TypeString,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - main - section title.",
						},
						"cardholder": {
							Type:        schema.TypeString,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - main - cardholder.",
						},
						"type": {
							Type:         schema.TypeString,
							Computed:     true,
							Optional:     true,
							ForceNew:     true,
							Description:  "Item credit card - main - type.",
							ValidateFunc: orEmpty(validation.StringInSlice([]string{"mc", "visa", "amex", "diners", "carteblanche", "discover", "jcb", "maestro", "visaelectron", "laser", "unionpay"}, false)),
						},
						"ccnum": {
							Type:        schema.TypeString,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - main - number.",
						},
						"cvv": {
							Type:        schema.TypeString,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: "Item credit card - main - cvv.",
						},
						"expiry": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"valid_from": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"field": sectionSchema.Schema["field"],
					},
				},
			},
			"contact_info": {
				Type:     schema.TypeList,
				Computed: false,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "Contact Information",
						},
						"bank": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - contact info - bank.",
							Elem:        labelValue("issuing bank"),
						},
						"phone_local": {
							Type:        schema.TypeSet,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - contact info - phone local.",
							Elem:        labelValue("phone (local)"),
						},
						"phone_toll_free": {
							Type:        schema.TypeSet,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - contact info - phone_toll_free.",
							Elem:        labelValue("phone (toll free)"),
						},
						"phone_intl": {
							Type:        schema.TypeSet,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - contact info - phone (intl).",
							Elem:        labelValue("phone (intl)"),
						},
						"website": {
							Type:        schema.TypeSet,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - contact info - website.",
							Elem:        labelValue("website"),
						},
						"field": sectionSchema.Schema["field"],
					},
				},
			},
			"details": {
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
							Default:  "Additional Details",
						},
						"pin": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     labelValue("PIN"),
						},
						"credit_limit": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     labelValue("credit limit"),
						},
						"cash_limit": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     labelValue("cash withdrawal limit"),
						},
						"interest": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							Elem:     labelValue("interest rate"),
						},
						"issuenumber": {
							Type:        schema.TypeSet,
							Computed:    false,
							Optional:    true,
							ForceNew:    true,
							Description: "Item credit card - additional info - issue number.",
							Elem:        labelValue("issue number"),
						},
						"field": sectionSchema.Schema["field"],
					},
				},
			},
			"section": {
				Type:        schema.TypeList,
				Computed:    false,
				Optional:    true,
				ForceNew:    true,
				Description: "Item credit card section.",
				Elem:        sectionSchema,
			},
		},
	}
}

func resourceItemCreditCardRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	log.Printf("[DEBUG] %v", v)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(CreditCardCategory) {
		return errors.New("Item is not from " + string(CreditCardCategory))
	}

	d.SetId(v.Uuid)
	d.Set("name", v.Overview.Title)
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)
	d.Set("notes", v.Details.Notes)
	err = parseSectionFromSchema(v.Details.Sections, d, []SectionGroup{
		SectionGroup{
			Name:     "main",
			Selector: "",
			Fields:   []string{"cardholder", "ccnum", "type", "cvv", "expiry", "validFrom"},
		},
		SectionGroup{
			Name:     "contact_info",
			Selector: "contactInfo",
			Fields:   []string{"bank", "phoneLocal", "phoneTollFree", "phoneIntl", "website"},
		},
		SectionGroup{
			Name:     "details",
			Selector: "details",
			Fields:   []string{"pin", "creditLimit", "cashLimit", "interest", "issuenumber"},
		},
	})
	if err != nil {
		return err
	}
	return d.Set("section", v.ProcessSections())
}

func resourceItemCreditCardCreate(d *schema.ResourceData, meta interface{}) error {
	main := d.Get("main").([]interface{})[0].(map[string]interface{})
	contactInfo := d.Get("contact_info").([]interface{})[0].(map[string]interface{})
	details := d.Get("details").([]interface{})[0].(map[string]interface{})
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
								Value: main["ccnum"].(string),
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
								Value: main["expiry"].(int),
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
					Section{
						Title: contactInfo["title"].(string),
						Name:  "contactInfo",
						Fields: append([]SectionField{
							SectionField{
								Type:  "string",
								Text:  contactInfo["bank"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: contactInfo["bank"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "bank",
								Inputs: map[string]string{
									"autocapitalization": "Words",
									"keyboard":           "Default",
								},
							},
							SectionField{
								Type:  "phone",
								Text:  contactInfo["phone_local"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: contactInfo["phone_local"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "phoneLocal",
								Inputs: map[string]string{
									"keyboard": "NamePhonePad",
								},
							},
							SectionField{
								Type:  "phone",
								Text:  contactInfo["phone_toll_free"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: contactInfo["phone_toll_free"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "phoneTollFree",
								Inputs: map[string]string{
									"keyboard": "NamePhonePad",
								},
							},
							SectionField{
								Type:  "phone",
								Text:  contactInfo["phone_intl"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: contactInfo["phone_intl"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "phoneIntl",
								Inputs: map[string]string{
									"keyboard": "NamePhonePad",
								},
							},
							SectionField{
								Type:  "URL",
								Text:  contactInfo["website"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: contactInfo["website"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "website",
							},
						}, ParseFields(main)...),
					},
					Section{
						Title: details["title"].(string),
						Name:  "details",
						Fields: append([]SectionField{
							SectionField{
								Type:  "concealed",
								Text:  details["pin"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: details["pin"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "pin",
								A: Annotation{
									generate: "off",
								},
								Inputs: map[string]string{
									"keyboard": "NumberPad",
								},
							},
							SectionField{
								Type:  "string",
								Text:  details["credit_limit"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: details["credit_limit"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "creditLimit",
								Inputs: map[string]string{
									"keyboard": "NumbersAndPunctuation",
								},
							},
							SectionField{
								Type:  "string",
								Text:  details["cash_limit"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: details["cash_limit"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "cashLimit",
								Inputs: map[string]string{
									"keyboard": "NumbersAndPunctuation",
								},
							},
							SectionField{
								Type:  "string",
								Text:  details["interest"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: details["interest"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "interest",
								Inputs: map[string]string{
									"keyboard": "NumbersAndPunctuation",
								},
							},
							SectionField{
								Type:  "string",
								Text:  details["issuenumber"].(*schema.Set).List()[0].(map[string]interface{})["label"].(string),
								Value: details["issuenumber"].(*schema.Set).List()[0].(map[string]interface{})["value"].(string),
								N:     "issuenumber",
								Inputs: map[string]string{
									"autocorrection": "no",
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
	d.SetId(item.Uuid)
	return nil
}
