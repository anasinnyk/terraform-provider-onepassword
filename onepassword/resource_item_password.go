package onepassword

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceItemPassword() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemPasswordRead,
		Create: resourceItemPasswordCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemPasswordRead(d, meta)
				return []*schema.ResourceData{d}, err
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
				Elem:     sectionSchema,
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: urlValidate,
			},
		},
	}
}

func resourceItemPasswordRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(PasswordCategory) {
		return errors.New("item is not from " + string(PasswordCategory))
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Overview.Title); err != nil {
		return err
	}
	if err := d.Set("url", v.Overview.URL); err != nil {
		return err
	}
	if err := d.Set("notes", v.Details.Notes); err != nil {
		return err
	}
	if err := d.Set("tags", v.Overview.Tags); err != nil {
		return err
	}
	if err := d.Set("vault", v.Vault); err != nil {
		return err
	}
	if err := d.Set("password", v.Details.Password); err != nil {
		return err
	}
	return d.Set("section", ProcessSections(v.Details.Sections))
}

func resourceItemPasswordCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(item.UUID)
	return nil
}
