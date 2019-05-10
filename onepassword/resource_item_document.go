package onepassword

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceItemDocument() *schema.Resource {
	return &schema.Resource{
		Read:   resourceItemDocumentRead,
		Create: resourceItemDocumentCreate,
		Delete: resourceItemDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceItemDocumentRead(d, meta)
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

func resourceItemDocumentRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	vaultID := d.Get("vault").(string)
	v, err := m.onePassClient.ReadItem(getID(d), vaultID)
	if err != nil {
		return err
	}
	if v.Template != Category2Template(DocumentCategory) {
		return errors.New("item is not from " + string(DocumentCategory))
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

	content, err := m.onePassClient.ReadDocument(v.UUID)
	if err != nil {
		return err
	}
	return d.Set("content", content)
}

func resourceItemDocumentCreate(d *schema.ResourceData, meta interface{}) error {
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
		return err
	}
	content, err := m.onePassClient.ReadDocument(item.UUID)
	if err != nil {
		return err
	}

	d.SetId(item.UUID)
	return d.Set("content", content)
}
