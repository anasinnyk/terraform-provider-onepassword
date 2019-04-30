package onepassword

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
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
			"onepassword_item_common":           resourceItemCommon(),
			"onepassword_item_software_license": resourceItemSoftwareLicense(),
			"onepassword_item_identity":         resourceItemIdentity(),
			"onepassword_item_password":         resourceItemPassword(),
			"onepassword_item_credit_card":      resourceItemCreditCard(),
			"onepassword_item_secure_note":      resourceItemSecureNote(),
			"onepassword_item_document":         resourceItemDocument(),
			"onepassword_item_login":            resourceItemLogin(),
			"onepassword_vault":                 resourceVault(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onepassword_item_common":           dataSourceItemCommon(),
			"onepassword_item_software_license": dataSourceItemSoftwareLicense(),
			"onepassword_item_identity":         dataSourceItemIdentity(),
			"onepassword_item_password":         dataSourceItemPassword(),
			"onepassword_item_credit_card":      dataSourceItemCreditCard(),
			"onepassword_item_secure_note":      dataSourceItemSecureNote(),
			"onepassword_item_document":         dataSourceItemDocument(),
			"onepassword_item_login":            dataSourceItemLogin(),
			"onepassword_vault":                 dataSourceVault(),
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

func unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}
		filenames = append(filenames, fpath)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return filenames, err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func installOPClient() (string, error) {
	version := "v0.5.5"
	if os.Getenv("OP_VERSION") != "" {
		version = os.Getenv("OP_VERSION")
	}
	binZip := fmt.Sprintf("/tmp/op_%s.zip", version)
	url := fmt.Sprintf(
		"https://cache.agilebits.com/dist/1P/op/pkg/%s/op_%s_%s_%s.zip",
		version,
		runtime.GOOS,
		runtime.GOARCH,
		version,
	)
	if _, err := os.Stat(binZip); os.IsNotExist(err) {
		resp, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		out, err := os.Create(binZip)
		if err != nil {
			return "", err
		}
		defer out.Close()
		if _, err = io.Copy(out, resp.Body); err != nil {
			return "", err
		}
		if _, err := unzip(binZip, "/tmp/terraform-provider-onepassword/"+version); err != nil {
			return "", err
		}
	}
	return "/tmp/terraform-provider-onepassword/" + version + "/op", nil
}

func (m *Meta) NewOnePassClient() (error, *OnePassClient) {
	bin, err := installOPClient()
	if err != nil {
		return err, nil
	}

	op := &OnePassClient{
		Email:     m.data.Get("email").(string),
		Password:  m.data.Get("password").(string),
		SecretKey: m.data.Get("secret_key").(string),
		Subdomain: m.data.Get("subdomain").(string),
		PathToOp:  bin,
		Session:   "",
	}
	if err := op.SignIn(); err != nil {
		return err, nil
	}
	return nil, op
}

func (o *OnePassClient) SignIn() error {
	cmd := exec.Command(o.PathToOp, "signin", o.Subdomain, o.Email, o.SecretKey, "--output=raw")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer stdin.Close()
		if _, err := io.WriteString(stdin, fmt.Sprintf("%s\n", o.Password)); err != nil {
			log.Println("[ERROR] ", err)
		}
	}()

	session, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	o.Session = string(session)
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
		err = fmt.Errorf("some error in command %v\nError: %s\nOutput: %s", args[:len(args)-1], err, res)
	}
	return err, res
}

func getResultId(r []byte) (error, string) {
	result := &Resource{}
	if err := json.Unmarshal(r, result); err != nil {
		return err, ""
	}
	return nil, result.UUID
}

type Resource struct {
	UUID string `json:"uuid"`
}

func getID(d *schema.ResourceData) string {
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
