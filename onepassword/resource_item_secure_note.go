package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceItemSecureNote() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceItemSecureNoteRead,
		CreateContext: resourceItemSecureNoteCreate,
		DeleteContext: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceItemSecureNoteRead(ctx, d, meta); err.HasError() {
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
			"section": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     sectionSchema(),
			},
		},
	}
}

func resourceItemSecureNoteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return diag.FromErr(err)
	}
	if v.Template != Category2Template(SecureNoteCategory) {
		return diag.FromErr(errors.New("item is not from " + string(SecureNoteCategory)))
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
	if err := d.Set("section", ProcessSections(v.Details.Sections)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceItemSecureNoteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(SecureNoteCategory),
		Details: Details{
			Notes:    d.Get("notes").(string),
			Sections: ParseSections(d),
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
