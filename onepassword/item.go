package onepassword

import (
	"log"
	"strings"
	"encoding/json"
	"github.com/kalaspuffar/base64url"
)

const ITEM_RESOURCE = "item"

type SectionFieldType string
type FieldType string

const (
	FieldPassword FieldType = "P"
	FieldText     FieldType = "T"
)

const (
	TypeAddress   SectionFieldType = "address"
	TypeString    SectionFieldType = "string"
	TypeURL       SectionFieldType = "URL"
	TypeEmail     SectionFieldType = "email"
	TypeDate      SectionFieldType = "date"
	TypeMonthYear SectionFieldType = "monthYear"
	TypeConcealed SectionFieldType = "concealed"
	TypePhone     SectionFieldType = "phone"
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
	Vault    string   `json:"vaultUuid"`
	Overview Overview `json:"overview"`
	Details  Details  `json:"details"`
}

type Details struct {
	Notes    string    `json:"notesPlain"`
	Fields   []Field   `json:"fields"`
	Sections []Section `json:"sections"`
}

type Section struct {
	Title  string         `json:"title"`
	Fields []SectionField `json:"fields"`
}

type SectionField struct {
	Type  SectionFieldType `json:"k"`
	Text  string           `json:"t"`
	Value interface{}      `json:"v"`
	N     string           `json:"n"`
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
	err, res := o.runCmd(ONE_PASSWORD_COMMAND_GET, ITEM_RESOURCE, id, "--vault", vaultId)
	if err != nil {
		return err, nil
	}
	if err = json.Unmarshal(res, item); err != nil {
		return err, nil
	}
	return nil, item
}

func (o *OnePassClient) CreateItem(v *Item, category string) (error, *Item) {
	details, err := json.Marshal(v.Details)
	if err != nil {
		return err, nil
	}
	log.Printf("[DEBUG] Store Items - %s", details)
	detailsHash := base64url.Encode([]byte(details))

	err, _ = o.runCmd(
		ONE_PASSWORD_COMMAND_CREATE,
		ITEM_RESOURCE,
		category,
		detailsHash,
		"--vault="+v.Vault,
		"--title="+v.Overview.Title,
		"--url="+v.Overview.Url,
		"--tags="+strings.Join(v.Overview.Tags, ","),
	)

	if err != nil {
		return err, nil
	}
	return nil, v
}

func (o *OnePassClient) DeleteItem(id string) error {
	return o.Delete(ITEM_RESOURCE, id)
}
