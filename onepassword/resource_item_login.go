package onepassword

import (
  "fmt"
  "net/url"
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/hashicorp/terraform/helper/validation"
)

func resourceItemLogin() *schema.Resource {
  return &schema.Resource{
    Read:   resourceItemLoginRead,
    Create: resourceItemLoginCreate,
    Delete: resourceItemLoginDelete,
    Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
        resourceItemLoginRead(d, meta)
				return []*schema.ResourceData{d}, nil
			},
		},
    Schema: map[string]*schema.Schema{
      "name": {
        Type:        schema.TypeString,
        Computed:    true,
        Required:    true,
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
      "tags": {
        Type:        schema.TypeList,
        Computed:    true,
        Optional:    true,
        Description: "Item login tags.",
        Elem:        &schema.Schema{Type: schema.TypeString},
      },
      "vault": {
        Type:             schema.TypeString,
        Computed:         true,
        Optional:         true,
        Description:      "Vault for item login.",
      },
      "section": {
        Type:             schema.TypeSet,
        Computed:         true,
        Optional:         true,
        Description:      "Item login section.",
        Elem:             &schema.Resource{
          Schema: map[string]*schema.Schema{
            "name": {
              Type:        schema.TypeString,
              Computed:    true,
              Optional:    true,
              Description: "Item login section name.",
            },
            "field":  {
              Type:        schema.TypeSet,
              Computed:    true,
              Optional:    true,
              Description: "Item login section field.",
              Elem:        &schema.Resource{
                Schema: map[string]*schema.Schema{
                  "name": {
                    Type:        schema.TypeString,
                    Computed:    true,
                    Optional:    true,
                    Description: "Item login section name.",
                  },
                  "type": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    Default:      "STRING",
                    Description:  "Item login section type.",
                    ValidateFunc: validation.StringInSlice(
                      []string{"STRING", "URL", "PHONE", "ADDRESS", "TOTP", "EMAIL", "DATE", "CONCEALED", "MONTH_YEAR"},
                      false,
                    ),
                  },
                  "string": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    Description:  "Item login section value for string.",
                  },
                  "url": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    ValidateFunc: urlValidate,
                    Description:  "Item login section value for url.",
                  },
                  "phone": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    Description:  "Item login section value for phone.",
                  },
                  "email": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    // ValidateFunc: emailValidate,
                    Description:  "Item login section value for email.",
                  },
                  "date": {
                    Type:         schema.TypeInt,
                    Computed:     true,
                    Optional:     true,
                    // ValidateFunc: dateValidate,
                    Description:  "Item login section value for date.",
                  },
                  "month_year": {
                    Type:         schema.TypeInt,
                    Computed:     true,
                    Optional:     true,
                    // ValidateFunc: monthYearValidate,
                    Description:  "Item login section value for month year.",
                  },
                  "totp": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    Sensitive:    true,
                    // ValidateFunc: totpValidate,
                    Description:  "Item login section value for totp.",
                  },
                  "concealed": {
                    Type:         schema.TypeString,
                    Computed:     true,
                    Optional:     true,
                    Sensitive:    true,
                    Description:  "Item login section value for password.",
                  },
                  "address": {
                    Type:         schema.TypeMap,
                    Computed:     true,
                    Optional:     true,
                    Description:  "Item login section value for address.",
                    Elem:         &schema.Resource{
                      Schema: map[string]*schema.Schema{
                        "country": {
                          Type:         schema.TypeString,
                          Computed:     true,
                          Optional:     true,
                          Description:  "Item login section value for address - country.",
                        },
                        "city": {
                          Type:         schema.TypeString,
                          Computed:     true,
                          Optional:     true,
                          Description:  "Item login section value for address - city.",
                        },
                        "region": {
                          Type:         schema.TypeString,
                          Computed:     true,
                          Optional:     true,
                          Description:  "Item login section value for address - region.",
                        },
                        "state": {
                          Type:         schema.TypeString,
                          Computed:     true,
                          Optional:     true,
                          Description:  "Item login section value for address - state.",
                        },
                        "street": {
                          Type:         schema.TypeString,
                          Computed:     true,
                          Optional:     true,
                          Description:  "Item login section value for address - street.",
                        },
                        "zip": {
                          Type:         schema.TypeString,
                          Computed:     true,
                          Optional:     true,
                          Description:  "Item login section value for address - zip.",
                        },
                      },
                    },
                  },
                },
              },
            },
          },
        },
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

func urlValidate(i interface{}, k string) (s []string, es []error) {
  v, ok := i.(string)
  if !ok {
    es = append(es, fmt.Errorf("expected type of %s to be string", k))
    return
  }
  _, err := url.ParseRequestURI(v)
  if err != nil {
    es = append(es, fmt.Errorf("%s is not an URL", v))
    return
  }
  return
}

func resourceItemLoginRead(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceItemLoginCreate(d *schema.ResourceData, meta interface{}) error {
  return nil
}

func resourceItemLoginDelete(d *schema.ResourceData, meta interface{}) error {
  return nil
}
