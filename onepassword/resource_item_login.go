package onepassword

import (
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceItemLogin() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemLoginRead,
		Create: resourceItemLoginCreate,
		Delete: resourceItemLoginDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemLoginRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login name.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login username.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Sensitive:   true,
				Description: "Item login password.",
			},
			"notes": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Item login notes.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Item login tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vault": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Vault for item login.",
			},
			"section": {
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Description: "Item login section.",
				Elem:        sectionSchema,
			},
			"url": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				Description:  "URL for item login.",
				ValidateFunc: urlValidate,
			},
		},
	}
}

func resourceItemLoginRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultId := d.Get("vault").(string)
	err, v := m.onePassClient.ReadItem(getId(d), vaultId)
	log.Printf("[DEBUG] %v", v)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(LoginCategory) {
		return errors.New("Item is not from " + string(LoginCategory))
	}

	d.SetId(v.Uuid)
	d.Set("name", v.Overview.Title)
	if err := d.Set("url", v.Overview.Url); err != nil {
		return err
	}
	d.Set("notes", v.Details.Notes)
	d.Set("tags", v.Overview.Tags)
	d.Set("vault", v.Vault)
	for _, field := range v.Details.Fields {
		if field.Name == "username" {
			d.Set("username", field.Value)
		}
		if field.Name == "password" {
			d.Set("password", field.Value)
		}
	}
	return d.Set("section", v.ProcessSections())
}

func resourceItemLoginCreate(d *schema.ResourceData, meta interface{}) error {
	item := &Item{
		Vault:    d.Get("vault").(string),
		Template: Category2Template(LoginCategory),
		Overview: Overview{
			Title: d.Get("name").(string),
			Url:   d.Get("url").(string),
			Tags:  ParseTags(d),
		},
		Details: Details{
			Notes: d.Get("notes").(string),
			Fields: []Field{
				Field{
					Name:        "username",
					Designation: "username",
					Value:       d.Get("username").(string),
					Type:        FieldText,
				},
				Field{
					Name:        "password",
					Designation: "password",
					Value:       d.Get("password").(string),
					Type:        FieldPassword,
				},
			},
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

func resourceItemLoginDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteItem(getId(d))
	if err == nil {
		d.SetId("")
	}
	return err
}
