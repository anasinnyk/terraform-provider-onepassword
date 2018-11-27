package onepassword

import (
  "log"
  "github.com/hashicorp/terraform/helper/schema"
)

func resourceVault() *schema.Resource {
  return &schema.Resource{
    Read:   resourceVaultRead,
    Create: resourceVaultCreate,
    Delete: resourceVaultDelete,
    Schema: map[string]*schema.Schema{
      "name": {
        Type:        schema.TypeString,
        ForceNew:    true,
        Required:    true,
        Description: "Vault name.",
      },
    },
  }
}

func resourceVaultRead(d *schema.ResourceData, meta interface{}) error {
  m := meta.(*Meta)
  err, v := m.onePassClient.ReadVault(getId(d));
  if err != nil {
    log.Print(err)
    return err
  }

  d.SetId(v.Uuid)
  d.Set("name", v.Name)
  return nil
}

func resourceVaultCreate(d *schema.ResourceData, meta interface{}) error {
  m := meta.(*Meta)
  err, _ := m.onePassClient.CreateVault(&Vault{
    Name: d.Get("name").(string),
  })
  if err != nil {
    log.Print(err)
    return err
  }
  return resourceVaultRead(d, meta)
}

func getId(d *schema.ResourceData) string {
  return d.Get("name").(string)
}

func resourceVaultDelete(d *schema.ResourceData, meta interface{}) error {
  m := meta.(*Meta)
  err := m.onePassClient.DeleteVault(getId(d))
  if err != nil {
    log.Print(err)
    return err
  }
  d.SetId("")
  return nil
}
