package onepassword

import (
  "io"
	"fmt"
  "log"
  "regexp"
  "strings"
  "os/exec"
  "encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set account email address",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set account password",
			},
      "secret_key": {
        Type:        schema.TypeString,
        Required:    true,
        Description: "Set account secret key",
      },
      "subdomain": {
        Type:        schema.TypeString,
        Optional:    true,
        Default:     "my",
        Description: "Set alternative subdomain for 1password. From [subdomain].1password.com",
      },
		},
		ResourcesMap: map[string]*schema.Resource{
      // "1password_item":     resourceItem(), TODO: InProgress
			// "1password_group":    resourceGroup(), TODO: check it in team account
			"1password_vault": resourceVault(),
      // "1password_user": resourceDocument(), TODO: InProgress
    },
    DataSourcesMap: map[string]*schema.Resource{
      "1password_vault": dataSourceVault(),
    },
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return NewMeta(d)
}

const ONE_PASSWORD_COMMAND_CREATE = "create"
const ONE_PASSWORD_COMMAND_DELETE = "delete"
const ONE_PASSWORD_COMMAND_GET    = "get"
const ONE_PASSWORD_COMMAND_UPDATE = "update"
const ONE_PASSWORD_COMMAND_ADD    = "add"

type OnePassClient struct {
  Password  string
  Email     string
  SecretKey string
  Subdomain string
  PathToOp  string
  Session   string
}

type Vault struct {
  Uuid        string
  Name        string
}

type Meta struct {
	data          *schema.ResourceData
  onePassClient *OnePassClient
}

func NewMeta(d *schema.ResourceData) (*Meta, error) {
	m := &Meta{data: d}
  err, client := m.NewOnePassClient()
  m.onePassClient = client
  return m, err
}

func (m *Meta) NewOnePassClient() (error, *OnePassClient) {
  op := &OnePassClient{
    Email:     m.data.Get("email").(string),
    Password:  m.data.Get("password").(string),
    SecretKey: m.data.Get("secret_key").(string),
    Subdomain: m.data.Get("subdomain").(string),
    PathToOp:  "/usr/local/bin/op",
    Session:   "",
  }
  if err := op.SignIn(); err != nil {
    log.Print(err)
    return err, nil
  }
  return nil, op
}

func (o *OnePassClient) SignIn() error {
  cmd := exec.Command(o.PathToOp, "signin", o.Subdomain, o.Email, o.SecretKey)
  stdin, err := cmd.StdinPipe()
  if err != nil {
    return err
  }
  go func() {
    defer stdin.Close()
    io.WriteString(stdin, fmt.Sprintf("%s\n", o.Password))
  }()

  out, err := cmd.CombinedOutput()
  if err != nil {
    log.Print("SIGNIN ERROR: ", err)
    return err
  }

  log.Print("SignIn Output: ", out)
  r := regexp.MustCompile(fmt.Sprintf("export OP_SESSION_%s=\"(.+)\"", strings.Replace(o.Subdomain, "-", "_", 1)))
  session := r.FindStringSubmatch(string(out))[1]
  if session == "" {
    return fmt.Errorf("Cannot parse session from output: %s", out)
  }
  o.Session = session
  return nil
}

func (o *OnePassClient) runCmd(args ...string) (error, []byte) {
  args = append(args, fmt.Sprintf("--session=%s", o.Session))
  cmd := exec.Command(o.PathToOp, args...)
  res, err := cmd.CombinedOutput()
  if err != nil {
    err = fmt.Errorf("Some error in command %v\nError: %s\nOutput: %s", args, err, res)
  }
  return err, res
}

func (o *OnePassClient) ReadVault(id string) (error, *Vault) {
  vault := &Vault{}
  err, res := o.runCmd(ONE_PASSWORD_COMMAND_GET, "vault", id)
  if err != nil {
    return err, nil
  }
  if err = json.Unmarshal(res, vault); err != nil {
    return err, nil
  }
  return nil, vault
}

func (o *OnePassClient) CreateVault(v *Vault) (error, *Vault) {
  args := []string{ONE_PASSWORD_COMMAND_CREATE, "vault", v.Name}
  err, res := o.runCmd(args...)
  if err != nil {
    return err, nil
  }
  if err = json.Unmarshal(res, v); err != nil {
    return err, nil
  }
  return nil, v
}

func (o *OnePassClient) DeleteVault(id string) error {
  err, _ := o.runCmd(ONE_PASSWORD_COMMAND_DELETE, "vault", id)
  return err
}
