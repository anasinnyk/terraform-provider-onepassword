package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceItemSoftwareLicense() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceItemSoftwareLicenseRead,
		CreateContext: resourceItemSoftwareLicenseCreate,
		DeleteContext: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceItemSoftwareLicenseRead(ctx, d, meta); err.HasError() {
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

func resourceItemSoftwareLicenseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return diag.FromErr(err)
	}
	if v.Template != Category2Template(SoftwareLicenseCategory) {
		return diag.FromErr(errors.New("item is not from " + string(SoftwareLicenseCategory)))
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
				"license_key": "reg_code",
			},
		},
	}); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceItemSoftwareLicenseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}
	d.SetId(item.UUID)
	return nil
}
