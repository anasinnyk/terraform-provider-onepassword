package onepassword

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/kalaspuffar/base64url"
	"log"
	"strings"
)

const ITEM_RESOURCE = "item"
const DOCUMENT_RESOURCE = "document"

type SectionFieldType string
type FieldType string
type Category string

const (
	FieldPassword FieldType = "P"
	FieldText     FieldType = "T"
)

const (
	LoginCategory                Category = "Login"
	IdentiryCategory             Category = "Identity"
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
	Uuid     string   `json:"uuid"`
	Template string   `json:"templateUuid"`
	Vault    string   `json:"vaultUuid"`
	Overview Overview `json:"overview"`
	Details  Details  `json:"details"`
}

type DocumentAttributes struct {
	FileName string `json:"fileName,omitempty"`
}

type Details struct {
	DocumentAttributes DocumentAttributes `json:"documentAttributes,omitempty"`
	Notes              string             `json:"notesPlain"`
	Fields             []Field            `json:"fields"`
	Sections           []Section          `json:"sections"`
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
	Fields   []string
}

type Annotation struct {
	generate        string
	guarded         string
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
	Url   string   `json:"url"`
	Tags  []string `json:"tags"`
}

func (o *OnePassClient) ReadItem(id string, vaultId string) (error, *Item) {
	item := &Item{}
	args := []string{
		ONE_PASSWORD_COMMAND_GET,
		ITEM_RESOURCE,
		id,
	}

	if vaultId != "" {
		args = append(args, fmt.Sprintf("--vault=%s", vaultId))
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
	case IdentiryCategory:
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
		return IdentiryCategory
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
	log.Printf("[DEBUG] Store Items - %s", details)
	detailsHash := base64url.Encode([]byte(details))
	template := Template2Category(v.Template)
	if template == UnknownCategory {
		return errors.New("Unknown template id " + v.Template)
	}

	args := []string{
		ONE_PASSWORD_COMMAND_CREATE,
		ITEM_RESOURCE,
		string(template),
		detailsHash,
		fmt.Sprintf("--title=%s", v.Overview.Title),
		fmt.Sprintf("--url=%s", v.Overview.Url),
		fmt.Sprintf("--tags=%s", strings.Join(v.Overview.Tags, ",")),
	}

	if v.Vault != "" {
		args = append(args, fmt.Sprintf("--vault=%s", v.Vault))
	}
	err, res := o.runCmd(args...)
	if err == nil {
		err, id := getResultId(res)
		if err == nil {
			v.Uuid = id
		}
	}
	return err
}

func (o *OnePassClient) ReadDocument(id string) (error, string) {
	err, content := o.runCmd(
		ONE_PASSWORD_COMMAND_GET,
		DOCUMENT_RESOURCE,
		id,
	)
	return err, string(content)
}

func (o *OnePassClient) CreateDocument(v *Item, filePath string) error {
	args := []string{
		ONE_PASSWORD_COMMAND_CREATE,
		DOCUMENT_RESOURCE,
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
			v.Uuid = id
		}
	}
	return err
}

func resourceItemDelete(d *schema.ResourceData, meta interface{}) error {
	m := meta.(*Meta)
	err := m.onePassClient.DeleteItem(getId(d))
	if err == nil {
		d.SetId("")
	}
	return err
}

func (o *OnePassClient) DeleteItem(id string) error {
	return o.Delete(ITEM_RESOURCE, id)
}

func (i *Item) ProcessSections() []map[string]interface{} {
	sections := make([]map[string]interface{}, 0, len(i.Details.Sections))
	for _, section := range i.Details.Sections {
		fields := make([]map[string]interface{}, 0, len(section.Fields))
		for _, field := range section.Fields {
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
		sections = append(sections, map[string]interface{}{
			"name":  section.Title,
			"field": fields,
		})
	}
	return sections
}

func parseSectionFromSchema(sections []Section, d *schema.ResourceData, groups []SectionGroup) error {
	for _, section := range sections {
		for _, group := range groups {
			if section.Name == group.Selector {
				var leftFields []SectionField
				src := map[string]interface{}{
					"title": section.Title,
				}
				for _, field := range section.Fields {
					for _, f := range group.Fields {
						if f == field.N {
							if field.A.guarded == "yes" {
								src[ToSnakeCase(f)] = field.Value
							} else {
								src[ToSnakeCase(f)] = map[string]interface{}{
									"value": field.Value,
									"label": field.Text,
								}
							}
						} else {
							leftFields = append(leftFields, field)
						}
					}
				}
				src["field"] = ParseFields(map[string]interface{}{
					"field": leftFields,
				})
				if err := d.Set(group.Name, src); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
