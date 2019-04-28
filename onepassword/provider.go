package onepassword

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"os/exec"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OP_EMAIL", nil),
				Description: "Set account email address",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OP_PASSWORD", nil),
				Description: "Set account password",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OP_SECRET_KEY", nil),
				Description: "Set account secret key",
			},
			"subdomain": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "my",
				DefaultFunc: schema.EnvDefaultFunc("OP_SUBDOMAIN", nil),
				Description: "Set alternative subdomain for 1password. From [subdomain].1password.com",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			// "op_item_wireless_router": resourceItemWirelessRouter(),
			// "op_item_software_license": resourceItemSoftwareLicense(),
			// "op_item_social_security_number": resourceItemSocialSecurityNumber(),
			// "op_item_server": resourceItemServer(),
			// "op_item_reward_program": resourceItemRewardProgram(),
			// "op_item_passport": resourceItemPassport(),
			// "op_item_outdoor_license": resourceItemOutdoorLicense(),
			// "op_item_membership": resourceItemMembership(),
			// "op_item_email_account": resourceItemEmailAccount(),
			// "op_item_driver_license": resourceItemDriverLicense(),
			// "op_item_database": resourceItemDatabase(),
			// "op_item_bank_account": resourceItemBankAccount(),
			// "op_item_identity":    resourceItemIdentity(),
			"op_item_credit_card": resourceItemCreditCard(),
			"op_item_secure_note": resourceItemSecureNote(),
			"op_item_document":    resourceItemDocument(),
			"op_item_login":       resourceItemLogin(),
			"op_vault":            resourceVault(),
			"op_group":            resourceGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"op_item_credit_card": dataSourceItemCreditCard(),
			"op_item_secure_note": dataSourceItemSecureNote(),
			"op_item_document":    dataSourceItemDocument(),
			"op_item_login":       dataSourceItemLogin(),
			"op_vault":            dataSourceVault(),
			"op_group":            dataSourceGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return NewMeta(d)
}

const ONE_PASSWORD_COMMAND_CREATE = "create"
const ONE_PASSWORD_COMMAND_DELETE = "delete"
const ONE_PASSWORD_COMMAND_UPDATE = "update"
const ONE_PASSWORD_COMMAND_GET = "get"
const ONE_PASSWORD_COMMAND_ADD = "add"

type OnePassClient struct {
	Password  string
	Email     string
	SecretKey string
	Subdomain string
	PathToOp  string
	Session   string
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
		log.Print("[ERROR] ", err)
		return err
	}

	log.Print("[DEBUG] SignIn Output: ", out)
	r := regexp.MustCompile(fmt.Sprintf("export OP_SESSION_%s=\"(.+)\"", strings.Replace(o.Subdomain, "-", "_", 1)))
	session := r.FindStringSubmatch(string(out))[1]
	if session == "" {
		return fmt.Errorf("Cannot parse session from output: %s", out)
	}
	o.Session = session
	return nil
}

var m sync.Mutex

func (o *OnePassClient) runCmd(args ...string) (error, []byte) {
	args = append(args, fmt.Sprintf("--session=%s", o.Session))
	m.Lock()
	cmd := exec.Command(o.PathToOp, args...)
	defer m.Unlock()
	res, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Some error in command %v\nError: %s\nOutput: %s", args, err, res)
	}
	return err, res
}

func getResultId(r []byte) (error, string) {
	result := &Resource{}
	if err := json.Unmarshal(r, result); err != nil {
		return err, ""
	}
	return nil, result.Uuid
}

type Resource struct {
	Uuid string `json:"uuid"`
}

func getId(d *schema.ResourceData) string {
	if d.Id() != "" {
		return d.Id()
	} else {
		return d.Get("name").(string)
	}
}

func (o *OnePassClient) Delete(resource string, id string) error {
	err, _ := o.runCmd(ONE_PASSWORD_COMMAND_DELETE, resource, id)
	return err
}
