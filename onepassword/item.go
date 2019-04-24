package onepassword

import (
	"encoding/json"
)

const ITEM_RESOURCE = "item"

type SectionFieldType string
type FieldType string

const (
	FieldPassword FieldType = "P"
	FieldText     FieldType = "T"
)

const (
	TypeAddress    SectionFieldType = "address"
	TypeString     SectionFieldType = "string"
	TypeURL        SectionFieldType = "URL"
	TypeEmail      SectionFieldType = "email"
	TypeDate       SectionFieldType = "date"
	TypeMounthYear SectionFieldType = "mounthYear"
	TypeConcealed  SectionFieldType = "concealed"
	TypePhone      SectionFieldType = "phone"
)

type Address struct {
	City    string
	Country string
	Region  string
	State   string
	Street  string
	zip     string
}

type Item struct {
	Uuid     string
	Vault    string `json:"vaultUuid"`
	Overview Overview
	Details  Details
}

type Details struct {
	Notes    string `json:"notesPlain"`
	Fields   []Field
	Sections []Section
}

type Section struct {
	Title  string
	Fields []SectionField
}

type SectionField struct {
	Type  SectionFieldType `json:"k"`
	Text  string           `json:"t"`
	Value interface{}      `json:"v"`
	N     string
}

type Field struct {
	Type  FieldType
	Name  string
	Value string
}

type Overview struct {
	Title string
	Url   string
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

func (o *OnePassClient) CreateItem(v *Item) (error, *Vault) {
	return nil, nil
}

func (o *OnePassClient) DeleteItem(id string) error {
	return o.Delete(ITEM_RESOURCE, id)
}
