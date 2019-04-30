package onepassword

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceItemSoftwareLicense() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemSoftwareLicenseRead,
		Create: resourceItemSoftwareLicenseCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemSoftwareLicenseRead(d, meta)
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
						"license_key": {
							Type:     schema.TypeString,
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

func resourceItemSoftwareLicenseRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(SoftwareLicenseCategory) {
		return errors.New("item is not from " + string(SoftwareLicenseCategory))
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
			Name:     "main",
			Selector: "",
			Fields: map[string]string{
				"license_key": "reg_code",
			},
		},
	})
}

func resourceItemSoftwareLicenseCreate(d *schema.ResourceData, meta interface{}) error {
	main := d.Get("main").([]interface{})[0].(map[string]interface{})
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(SoftwareLicenseCategory),
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
								Text:  "license key",
								Value: main["license_key"].(string),
								N:     "reg_code",
								A: Annotation{
									guarded:   "yes",
									multiline: "yes",
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
