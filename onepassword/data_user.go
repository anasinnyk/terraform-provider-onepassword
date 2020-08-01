package onepassword

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
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

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	v, err := m.onePassClient.ReadUser(getIDEmail(d))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(v.UUID)
	if err := d.Set("email", v.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("firstname", v.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("lastname", v.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", v.State); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
