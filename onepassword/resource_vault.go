package onepassword

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVault() *schema.Resource {
	return &schema.Resource{
		Read:   resourceVaultRead,
		Create: resourceVaultCreate,
		Delete: resourceVaultDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceVaultRead(d, meta)
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

func resourceVaultRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err, v := m.onePassClient.ReadVault(getId(d))
	if err != nil {
		return err
	}

	d.SetId(v.Uuid)
	return d.Set("name", v.Name)
}

func resourceVaultCreate(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err, _ := m.onePassClient.CreateVault(&Vault{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return err
	}
	return resourceVaultRead(d, meta)
}

func resourceVaultDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteVault(getId(d))
	if err == nil {
		d.SetId("")
	}
	return err
}
