package onepassword

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var version string = "1.4.0"

func Provider() *schema.Provider {
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
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					if v := os.Getenv("OP_SUBDOMAIN"); v != "" {
						return v, nil
					}
					return "my", nil
				},
				Description: "Set alternative subdomain for 1password. From [subdomain].1password.com",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"onepassword_group":                 resourceGroup(),
			"onepassword_group_member":          resourceGroupMember(),
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
			"onepassword_group":                 dataSourceGroup(),
			"onepassword_user":                  dataSourceUser(),
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
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return NewMeta(d)
}

const (
	opPasswordAdd    = "add"
	opPasswordCreate = "create"
	opPasswordEdit   = "edit"
	opPasswordDelete = "delete"
	opPasswordGet    = "get"
	opPasswordList   = "list"
	opPasswordRemove = "remove"
)

type OnePassClient struct {
	Password    string
	Email       string
	SecretKey   string
	Subdomain   string
	PathToOp    string
	Session     string
	execCommand func(string, ...string) *exec.Cmd // Can be overridden for mocking purposes
	mutex       *sync.Mutex
}

type Meta struct {
	data          *schema.ResourceData
	onePassClient *OnePassClient
}

func NewMeta(d *schema.ResourceData) (*Meta, diag.Diagnostics) {
	m := &Meta{data: d}
	client, err := m.NewOnePassClient()
	if err != nil {
		return m, diag.FromErr(err)
	}
	m.onePassClient = client
	return m, nil
}

func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		traversableCheck := strings.Split(f.Name, "..")
		fpath := filepath.Join(dest, traversableCheck[len(traversableCheck)-1])
		if err != nil {
			return err
		}
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}

func findExistingOPClient() (string, error) {
	o, err := exec.Command("op", "--version").Output()

	if err != nil {
		return "", fmt.Errorf("Trouble calling: op\nOutput: %s", o)
	}

	c, err := semver.NewConstraint(">= " + version)
	if err != nil {
		return "", err
	}

	v, err := semver.NewVersion(strings.TrimSuffix(string(o), "\n"))
	if err != nil {
		return "", fmt.Errorf("[%s]", string(o))
	}

	if c.Check(v) {
		return "op", nil
	}

	return "", fmt.Errorf("op version needs to be equal or greater than: %s", version)
}

func installOPClient() (string, error) {
	if os.Getenv("OP_VERSION") != "" {
		semVer, err := semver.NewVersion(os.Getenv("OP_VERSION"))
		if err != nil {
			return "", err
		}
		version = semVer.String()
	}
	if runtime.GOOS == "darwin" {
		return "", fmt.Errorf("Unable to automatically install v%s of the op client. Please install manually from https://app-updates.agilebits.com/product_history/CLI", version)
	}

	binZip := fmt.Sprintf("/tmp/op_%s.zip", version)
	if _, err := os.Stat(binZip); os.IsNotExist(err) {
		resp, err := http.Get(fmt.Sprintf(
			"https://cache.agilebits.com/dist/1P/op/pkg/v%s/op_%s_%s_v%s.zip",
			version,
			runtime.GOOS,
			runtime.GOARCH,
			version,
		))
		if err != nil {
			return "", fmt.Errorf("Could not retrieve zipped op release: %w", err)
		}
		defer resp.Body.Close()

		out, err := os.Create(binZip)
		if err != nil {
			return "", fmt.Errorf("Could not create temp file for op client: %w", err)
		}
		defer out.Close()
		if _, err = io.Copy(out, resp.Body); err != nil {
			return "", fmt.Errorf("Could not copy zip contents to temp file for op client: %w", err)
		}
		if err := unzip(binZip, "/tmp/terraform-provider-onepassword/"+version); err != nil {
			return "", fmt.Errorf("Could not unzip temp file for op client: %w", err)
		}
	}
	return "/tmp/terraform-provider-onepassword/" + version + "/op", nil
}

func (m *Meta) NewOnePassClient() (*OnePassClient, error) {
	bin, err := findExistingOPClient()
	if err != nil {
		bin, err = installOPClient()
		if err != nil {
			return nil, err
		}
	}

	subdomain := m.data.Get("subdomain").(string)
	email := m.data.Get("email").(string)
	password := m.data.Get("password").(string)
	secretKey := m.data.Get("secret_key").(string)
	session := ""

	if email == "" || password == "" || secretKey == "" {
		email = ""
		password = ""
		secretKey = ""

		var sessionKeyName string
		if strings.Contains(subdomain, "-") {
			sessionKeyName = "OP_SESSION_" + strings.ReplaceAll(subdomain, "-", "_")
		} else {
			sessionKeyName = "OP_SESSION_" + subdomain
		}
		session = os.Getenv(sessionKeyName)

		if session == "" {
			return nil, fmt.Errorf("email, password or secret_key is empty and environment variable %s is not set",
				sessionKeyName)
		}
	}

	op := &OnePassClient{
		Email:       email,
		Password:    password,
		SecretKey:   secretKey,
		Subdomain:   subdomain,
		PathToOp:    bin,
		Session:     session,
		execCommand: exec.Command,
		mutex:       &sync.Mutex{},
	}

	if session != "" {
		return op, nil
	}
	if err := op.SignIn(); err != nil {
		return nil, err
	}
	return op, nil
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
		return errors.New(fmt.Sprintf("Cannot signin: %s\nExit code: %s", string(session), err))
	}

	o.Session = string(session)
	return nil
}

func (o *OnePassClient) runCmd(args ...string) ([]byte, error) {
	args = append(args, fmt.Sprintf("--session=%s", strings.Trim(o.Session, "\n")))
	o.mutex.Lock()
	cmd := o.execCommand(o.PathToOp, args...)
	defer o.mutex.Unlock()
	res, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("some error in command %v\nError: %s\nOutput: %s", args[:len(args)-1], err, res)
	}
	return res, err
}

func getResultID(r []byte) (string, error) {
	result := &Resource{}
	if err := json.Unmarshal(r, result); err != nil {
		return "", err
	}
	return result.UUID, nil
}

type Resource struct {
	UUID string `json:"uuid"`
}

func getID(d *schema.ResourceData) string {
	if d.Id() != "" {
		return d.Id()
	}
	return d.Get("name").(string)
}

func getIDEmail(d *schema.ResourceData) string {
	if d.Id() != "" {
		return d.Id()
	}
	return d.Get("email").(string)
}

func (o *OnePassClient) Delete(resource string, id string) error {
	_, err := o.runCmd(opPasswordDelete, resource, id)
	return err
}
