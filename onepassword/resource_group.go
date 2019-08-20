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
	v, err := m.onePassClient.ReadGroup(getID(d))
	if err != nil {
		return err
	} else if v.State == GroupStateDeleted {
		d.SetId("")
		return nil
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Name); err != nil {
		return err
	}

	return d.Set("state", v.State)
}

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	_, err := m.onePassClient.CreateGroup(&Group{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteGroup(getID(d))
	if err == nil {
		d.SetId("")
	}
	return err
}
