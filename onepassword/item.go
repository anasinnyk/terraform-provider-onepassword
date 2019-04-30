package onepassword

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kalaspuffar/base64url"
	"strings"
)

const ItemResource = "item"
const DocumentResource = "document"

type SectionFieldType string
type FieldType string
type Category string

const (
	FieldPassword FieldType = "P"
	FieldText     FieldType = "T"
)

const (
	LoginCategory                Category = "Login"
	IdentityCategory             Category = "Identity"
	DatabaseCategory             Category = "Database"
	MembershipCategory           Category = "Membership"
	WirelessRouterCategory       Category = "Wireless Router"
	SecureNoteCategory           Category = "Secure Note"
	SoftwareLicenseCategory      Category = "Software License"
	CreditCardCategory           Category = "Credit Card"
	DriverLicenseCategory        Category = "Driver License"
	OutdoorLicenseCategory       Category = "Outdoor License"
	PassportCategory             Category = "Passport"
	EmailAccountCategory         Category = "Email Account"
	PasswordCategory             Category = "Password"
	RewardProgramCategory        Category = "Reward Program"
	SocialSecurityNumberCategory Category = "Social Security Number"
	BankAccountCategory          Category = "Bank Account"
	DocumentCategory             Category = "Document"
	ServerCategory               Category = "Server"
	UnknownCategory              Category = "UNKNOWN"
)

const (
	TypeSex       SectionFieldType = "menu"
	TypeCard      SectionFieldType = "cctype"
	TypeAddress   SectionFieldType = "address"
	TypeString    SectionFieldType = "string"
	TypeURL       SectionFieldType = "URL"
	TypeEmail     SectionFieldType = "email"
	TypeDate      SectionFieldType = "date"
	TypeMonthYear SectionFieldType = "monthYear"
	TypeConcealed SectionFieldType = "concealed"
	TypePhone     SectionFieldType = "phone"
	TypeReference SectionFieldType = "reference"
)

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Region  string `json:"region"`
	State   string `json:"state"`
	Street  string `json:"street"`
	Zip     string `json:"zip"`
}

type Item struct {
	UUID     string   `json:"uuid"`
	Template string   `json:"templateUUID"`
	Vault    string   `json:"vaultUUID"`
	Overview Overview `json:"overview"`
	Details  Details  `json:"details"`
}

type Details struct {
	Notes    string    `json:"notesPlain"`
	Password string    `json:"password"`
	Fields   []Field   `json:"fields"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Name   string         `json:"name"`
	Title  string         `json:"title"`
	Fields []SectionField `json:"fields"`
}

type SectionField struct {
	Type   SectionFieldType  `json:"k"`
	Text   string            `json:"t"`
	Value  interface{}       `json:"v"`
	N      string            `json:"n"`
	A      Annotation        `json:"a"`
	Inputs map[string]string `json:"inputTraits"`
}

type SectionGroup struct {
	Selector string
	Name     string
	Fields   map[string]string
}

type Annotation struct {
	generate        string
	guarded         string
	multiline       string
	clipboardFilter string
}

type Field struct {
	Type        FieldType `json:"type"`
	Designation string    `json:"designation"`
	Name        string    `json:"name"`
	Value       string    `json:"value"`
}

type Overview struct {
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags"`
}

func (o *OnePassClient) ReadItem(id string, vaultID string) (error, *Item) {
	item := &Item{}
	args := []string{
		ONE_PASSWORD_COMMAND_GET,
		ItemResource,
		id,
	}

	if vaultID != "" {
		args = append(args, fmt.Sprintf("--vault=%s", vaultID))
	}
	err, res := o.runCmd(args...)
	if err != nil {
		return err, nil
	}
	if err = json.Unmarshal(res, item); err != nil {
		return err, nil
	}
	return nil, item
}

func Category2Template(c Category) string {
	switch c {
	case LoginCategory:
		return "001"
	case IdentityCategory:
		return "004"
	case PasswordCategory:
		return "005"
	case PassportCategory:
		return "106"
	case DatabaseCategory:
		return "102"
	case ServerCategory:
		return "010"
	case DriverLicenseCategory:
		return "103"
	case OutdoorLicenseCategory:
		return "104"
	case SoftwareLicenseCategory:
		return "100"
	case EmailAccountCategory:
		return "111"
	case RewardProgramCategory:
		return "107"
	case WirelessRouterCategory:
		return "109"
	case DocumentCategory:
		return "006"
	case BankAccountCategory:
		return "101"
	case SocialSecurityNumberCategory:
		return "108"
	case CreditCardCategory:
		return "002"
	case SecureNoteCategory:
		return "003"
	case MembershipCategory:
		return "105"
	default:
		return "000"
	}
}

func Template2Category(t string) Category {
	switch t {
	case "001":
		return LoginCategory
	case "004":
		return IdentityCategory
	case "005":
		return PasswordCategory
	case "106":
		return PassportCategory
	case "102":
		return DatabaseCategory
	case "010":
		return ServerCategory
	case "103":
		return DriverLicenseCategory
	case "104":
		return OutdoorLicenseCategory
	case "100":
		return SoftwareLicenseCategory
	case "111":
		return EmailAccountCategory
	case "107":
		return RewardProgramCategory
	case "109":
		return WirelessRouterCategory
	case "006":
		return DocumentCategory
	case "101":
		return BankAccountCategory
	case "108":
		return SocialSecurityNumberCategory
	case "002":
		return CreditCardCategory
	case "003":
		return SecureNoteCategory
	case "105":
		return MembershipCategory
	default:
		return UnknownCategory
	}
}

func (o *OnePassClient) CreateItem(v *Item) error {
	details, err := json.Marshal(v.Details)
	if err != nil {
		return err
	}
	detailsHash := base64url.Encode([]byte(details))
	template := Template2Category(v.Template)
	if template == UnknownCategory {
		return errors.New("unknown template id " + v.Template)
	}

	args := []string{
		ONE_PASSWORD_COMMAND_CREATE,
		ItemResource,
		string(template),
		detailsHash,
		fmt.Sprintf("--title=%s", v.Overview.Title),
		fmt.Sprintf("--url=%s", v.Overview.URL),
		fmt.Sprintf("--tags=%s", strings.Join(v.Overview.Tags, ",")),
	}

	if v.Vault != "" {
		args = append(args, fmt.Sprintf("--vault=%s", v.Vault))
	}
	err, res := o.runCmd(args...)
	if err == nil {
		err, id := getResultId(res)
		if err == nil {
			v.UUID = id
		}
	}
	return err
}

func (o *OnePassClient) ReadDocument(id string) (error, string) {
	err, content := o.runCmd(
		ONE_PASSWORD_COMMAND_GET,
		DocumentResource,
		id,
	)
	return err, string(content)
}

func (o *OnePassClient) CreateDocument(v *Item, filePath string) error {
	args := []string{
		ONE_PASSWORD_COMMAND_CREATE,
		DocumentResource,
		filePath,
		fmt.Sprintf("--title=%s", v.Overview.Title),
		fmt.Sprintf("--tags=%s", strings.Join(v.Overview.Tags, ",")),
	}

	if v.Vault != "" {
		args = append(args, fmt.Sprintf("--vault=%s", v.Vault))
	}

	err, res := o.runCmd(args...)
	if err == nil {
		err, id := getResultId(res)
		if err == nil {
			v.UUID = id
		}
	}
	return err
}

func resourceItemDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteItem(getID(d))
	if err == nil {
		d.SetId("")
	}
	return err
}

func (o *OnePassClient) DeleteItem(id string) error {
	return o.Delete(ItemResource, id)
}

func ProcessField(srcFields []SectionField) []map[string]interface{} {
	fields := make([]map[string]interface{}, 0, len(srcFields))
	for _, field := range srcFields {
		f := map[string]interface{}{
			"name": field.Text,
		}
		var key string
		switch field.Type {
		case TypeSex:
			key = "sex"
		case TypeURL:
			key = "url"
		case TypeMonthYear:
			key = "month_year"
		case TypeCard:
			key = "card_type"
		case TypeConcealed:
			if strings.HasPrefix(field.N, "TOTP_") {
				key = "totp"
			} else {
				key = "concealed"
			}
		default:
			key = string(field.Type)
		}
		f[key] = field.Value
		fields = append(fields, f)
	}
	return fields
}

func ProcessSections(srcSections []Section) []map[string]interface{} {
	sections := make([]map[string]interface{}, 0, len(srcSections))
	for _, section := range srcSections {
		sections = append(sections, map[string]interface{}{
			"name":  section.Title,
			"field": ProcessField(section.Fields),
		})
	}
	return sections
}

func parseSectionFromSchema(sections []Section, d *schema.ResourceData, groups []SectionGroup) error {
	leftSections := []Section{}
	for _, section := range sections {
		var use bool
		for _, group := range groups {
			if section.Name == group.Selector {
				use = true
				var leftFields []SectionField
				src := map[string]interface{}{
					"title": section.Title,
				}
				for _, field := range section.Fields {
					found := false
					for k, f := range group.Fields {
						if f == field.N {
							src[k] = field.Value
							found = true
							continue
						}
					}
					if !found {
						leftFields = append(leftFields, field)
					}
				}
				src["field"] = ProcessField(leftFields)
				if err := d.Set(group.Name, []interface{}{src}); err != nil {
					return err
				}
			}
		}
		if !use {
			leftSections = append(leftSections, section)
		}
	}
	return d.Set("section", ProcessSections(leftSections))
}
