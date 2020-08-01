package onepassword

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceGroupRead,
		CreateContext: resourceGroupCreate,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if err := resourceGroupRead(ctx, d, meta); err.HasError() {
					return []*schema.ResourceData{d}, errors.New(err[0].Summary)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	v, err := m.onePassClient.ReadGroup(getID(d))
	if err != nil {
		return diag.FromErr(err)
	} else if v.State == GroupStateDeleted {
		d.SetId("")
		return nil
	}

	d.SetId(v.UUID)
	if err := d.Set("name", v.Name); err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("state", v.State)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	_, err := m.onePassClient.CreateGroup(&Group{
		Name: d.Get("name").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteGroup(getID(d))
	if err == nil {
		d.SetId("")
	}
	return diag.FromErr(err)
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)

	g := &Group{
		Name: d.Get("name").(string),
	}

	if err := m.onePassClient.UpdateGroup(getID(d), g); err != nil {
		return diag.FromErr(err)
	}
	return resourceGroupRead(ctx, d, meta)
}
