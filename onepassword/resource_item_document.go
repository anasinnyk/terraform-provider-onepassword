package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceItemDocument() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceItemDocumentRead,
		CreateContext: resourceItemDocumentCreate,
		DeleteContext: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceItemDocumentRead(ctx, d, meta); err.HasError() {
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
				ForceNew: true,
				Optional: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func resourceItemDocumentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return diag.FromErr(err)
	}
	if v.Template != Category2Template(DocumentCategory) {
		return diag.FromErr(errors.New("item is not from " + string(DocumentCategory)))
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

	content, err := m.onePassClient.ReadDocument(v.UUID)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("content", content); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceItemDocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(DocumentCategory),
		Overview: Overview{
			Title: d.Get("name").(string),
			Tags:  ParseTags(d),
		},
	}
	m := meta.(*Meta)
	err := m.onePassClient.CreateDocument(item, d.Get("file_path").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	content, err := m.onePassClient.ReadDocument(item.UUID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(item.UUID)
	if err := d.Set("content", content); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
