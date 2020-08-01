package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceItemCommon() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceItemCommonRead,
		CreateContext: resourceItemCommonCreate,
		DeleteContext: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceItemCommonRead(ctx, d, meta); err.HasError() {
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
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(DatabaseCategory),
					string(MembershipCategory),
					string(WirelessRouterCategory),
					string(DriverLicenseCategory),
					string(OutdoorLicenseCategory),
					string(PassportCategory),
					string(EmailAccountCategory),
					string(RewardProgramCategory),
					string(SocialSecurityNumberCategory),
					string(BankAccountCategory),
					string(ServerCategory),
				}, false),
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
		},
	}
}

func resourceItemCommonRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Overview.Title); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("notes", v.Details.Notes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", v.Overview.Tags); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vault", v.Vault); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("template", string(Template2Category(v.Template))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("section", ProcessSections(v.Details.Sections)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceItemCommonCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(Category(d.Get("template").(string))),
		Overview: Overview{
			Title: d.Get("name").(string),
			Tags:  ParseTags(d),
		},
		Details: Details{
			Notes:    d.Get("notes").(string),
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
