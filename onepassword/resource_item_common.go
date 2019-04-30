package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceItemCommon() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemCommonRead,
		Create: resourceItemCommonCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemCommonRead(d, meta)
				return []*schema.ResourceData{d}, err
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
				Elem:     sectionSchema,
			},
		},
	}
}

func resourceItemCommonRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	if err != nil {
		return err
	}

	d.SetId(v.Uuid)
	if err := d.Set("name", v.Overview.Title); err != nil {
		return err
	}
	if err := d.Set("notes", v.Details.Notes); err != nil {
		return err
	}
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)
	d.Set("template", string(Template2Category(v.Template)))
	return d.Set("section", ProcessSections(v.Details.Sections))
}

func resourceItemCommonCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(item.Uuid)
	return nil
}
