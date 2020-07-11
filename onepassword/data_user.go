package onepassword

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"firstname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lastname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	v, err := m.onePassClient.ReadUser(getIDEmail(d))
	if err != nil {
		return err
	}

	d.SetId(v.UUID)
	if err := d.Set("email", v.Email); err != nil {
		return err
	}
	if err := d.Set("firstname", v.FirstName); err != nil {
		return err
	}
	if err := d.Set("lastname", v.LastName); err != nil {
		return err
	}

	return d.Set("state", v.State)
}
