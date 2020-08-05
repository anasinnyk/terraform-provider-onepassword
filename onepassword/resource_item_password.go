package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceItemPassword() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceItemPasswordRead,
		CreateContext: resourceItemPasswordCreate,
		DeleteContext: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceItemPasswordRead(ctx, d, meta); err.HasError() {
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
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  true,
			},
			"notes": {
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
			"section": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     sectionSchema(),
			},
			"url": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateDiagFunc: urlValidateDiag(),
			},
		},
	}
}

func resourceItemPasswordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		diag.FromErr(err)
	}
	if v.Template != Category2Template(PasswordCategory) {
		diag.FromErr(errors.New("item is not from " + string(PasswordCategory)))
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Overview.Title); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("url", v.Overview.URL); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("notes", v.Details.Notes); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("tags", v.Overview.Tags); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("vault", v.Vault); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("password", v.Details.Password); err != nil {
		diag.FromErr(err)
	}
	if err := d.Set("section", ProcessSections(v.Details.Sections)); err != nil {
		diag.FromErr(err)
	}
	return nil
}

func resourceItemPasswordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(PasswordCategory),
		Overview: Overview{
			Title: d.Get("name").(string),
			URL:   d.Get("url").(string),
			Tags:  ParseTags(d),
		},
		Details: Details{
			Notes:    d.Get("notes").(string),
			Password: d.Get("password").(string),
			Sections: ParseSections(d),
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
