package onepassword

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceItemIdentity() *schema.Resource {
	addressSchema := sectionSchema().Schema["field"].Elem.(*schema.Resource).Schema["address"]
	addressSchema.ConflictsWith = []string{}

	return &schema.Resource{
		Read:   resourceItemIdentityRead,
		Create: resourceItemIdentityCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemIdentityRead(d, meta)
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
			"identification": {
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
							Default:  "Identification",
						},
						"firstname": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"initial": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"lastname": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"sex": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: orEmpty(validation.StringInSlice([]string{"male", "female"}, false)),
						},
						"birth_date": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"occupation": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"company": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"department": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"job_title": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"field": sectionSchema().Schema["field"],
					},
				},
			},
			"address": {
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
							Default:  "Address",
						},
						"address": addressSchema,
						"default_phone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"home_phone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"cell_phone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"business_phone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"field": sectionSchema().Schema["field"],
					},
				},
			},
			"internet": {
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
							Default:  "Internet Details",
						},
						"username": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"email": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: emailValidate,
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

func resourceItemIdentityRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(IdentityCategory) {
		return errors.New("item is not from " + string(IdentityCategory))
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
		{
			Name:     "identification",
			Selector: "name",
			Fields: map[string]string{
				"firstname":  "firstname",
				"initial":    "initial",
				"lastname":   "lastname",
				"sex":        "sex",
				"birth_date": "birthdate",
				"occupation": "occupation",
				"company":    "company",
				"department": "department",
				"job_title":  "jobtitle",
			},
		},
		{
			Name:     "address",
			Selector: "address",
			Fields: map[string]string{
				"address":        "address",
				"default_phone":  "defphone",
				"home_phone":     "homephone",
				"cell_phone":     "cellphone",
				"business_phone": "busphone",
			},
		},
		{
			Name:     "internet",
			Selector: "internet",
			Fields: map[string]string{
				"username": "username",
				"email":    "email",
			},
		},
	})
}

func resourceItemIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	main := d.Get("identification").([]interface{})[0].(map[string]interface{})
	address := d.Get("address").([]interface{})[0].(map[string]interface{})
	internet := d.Get("internet").([]interface{})[0].(map[string]interface{})
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(IdentityCategory),
		Details: Details{
			Notes: d.Get("notes").(string),
			Sections: append(
				[]Section{
					{
						Title: main["title"].(string),
						Name:  "name",
						Fields: append([]SectionField{
							{
								Type:  "string",
								Text:  "firstname",
								Value: main["firstname"].(string),
								N:     "firstname",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
							{
								Type:  "string",
								Text:  "initial",
								Value: main["initial"].(string),
								N:     "initial",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
							{
								Type:  "string",
								Text:  "lastname",
								Value: main["lastname"].(string),
								N:     "lastname",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
							{
								Type:  "menu",
								Text:  "sex",
								Value: main["sex"].(string),
								N:     "sex",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
								Type:  "date",
								Text:  "birth date",
								Value: main["birth_date"].(int),
								N:     "birthdate",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
								Type:  "string",
								Text:  "occupation",
								Value: main["occupation"].(string),
								N:     "occupation",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
							{
								Type:  "string",
								Text:  "company",
								Value: main["company"].(string),
								N:     "company",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
							{
								Type:  "string",
								Text:  "department",
								Value: main["department"].(string),
								N:     "department",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
							{
								Type:  "string",
								Text:  "job title",
								Value: main["job_title"].(string),
								N:     "jobtitle",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Words",
								},
							},
						}, ParseFields(main)...),
					},
					{
						Title: address["title"].(string),
						Name:  "address",
						Fields: append([]SectionField{
							{
								Type:  "address",
								Text:  "address",
								Value: address["address"].(map[string]interface{}),
								N:     "address",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Sentences",
								},
							},
							{
								Type:  "phone",
								Text:  "default phone",
								Value: address["default_phone"].(string),
								N:     "defphone",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
								Type:  "phone",
								Text:  "home",
								Value: address["home_phone"].(string),
								N:     "homephone",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
								Type:  "phone",
								Text:  "cell",
								Value: address["cell_phone"].(string),
								N:     "cellphone",
								A: Annotation{
									guarded: "yes",
								},
							},
							{
								Type:  "phone",
								Text:  "business",
								Value: address["business_phone"].(string),
								N:     "busphone",
								A: Annotation{
									guarded: "yes",
								},
							},
						}, ParseFields(address)...),
					},
					{
						Title: internet["title"].(string),
						Name:  "internet",
						Fields: append([]SectionField{
							{
								Type:  "string",
								Text:  "username",
								Value: internet["username"].(string),
								N:     "username",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"autocapitalization": "Sentences",
								},
							},
							{
								Type:  "string",
								Text:  "email",
								Value: internet["email"].(string),
								N:     "email",
								A: Annotation{
									guarded: "yes",
								},
								Inputs: map[string]string{
									"keyboard": "EmailAddress",
								},
							},
						}, ParseFields(internet)...),
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
