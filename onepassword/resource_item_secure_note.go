package onepassword

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceItemSecureNote() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemSecureNoteRead,
		Create: resourceItemSecureNoteCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemSecureNoteRead(d, meta)
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
			"section": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     sectionSchema,
			},
		},
	}
}

func resourceItemSecureNoteRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	log.Printf("[DEBUG] %v", v)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(SecureNoteCategory) {
		return errors.New("Item is not from " + string(SecureNoteCategory))
	}

	d.SetId(v.Uuid)
	d.Set("name", v.Overview.Title)
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)
	d.Set("notes", v.Details.Notes)
	return d.Set("section", ProcessSections(v.Details.Sections))
}

func resourceItemSecureNoteCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	d.SetId(item.Uuid)
	return nil
}