package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Read:   resourceGroupRead,
		Create: resourceGroupCreate,
		Delete: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceGroupRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err, v := m.onePassClient.ReadGroup(getId(d))
	if err != nil {
		return err
	}

	d.SetId(v.Uuid)
	return d.Set("name", v.Name)
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err, _ := m.onePassClient.CreateGroup(&Group{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteGroup(getId(d))
	if err == nil {
		d.SetId("")
	}
	return err
}
