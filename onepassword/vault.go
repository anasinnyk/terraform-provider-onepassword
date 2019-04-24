package onepassword

import (
	"encoding/json"
)

const VAULT_RESOURCE = "vault"

type Vault struct {
	Uuid string
	Name string
}

func (o *OnePassClient) ReadVault(id string) (error, *Vault) {
	vault := &Vault{}
	err, res := o.runCmd(ONE_PASSWORD_COMMAND_GET, VAULT_RESOURCE, id)
	if err != nil {
		return err, nil
	}
	if err = json.Unmarshal(res, vault); err != nil {
		return err, nil
	}
	return nil, vault
}

func (o *OnePassClient) CreateVault(v *Vault) (error, *Vault) {
	args := []string{ONE_PASSWORD_COMMAND_CREATE, VAULT_RESOURCE, v.Name}
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
	err, _ := o.runCmd(ONE_PASSWORD_COMMAND_DELETE, VAULT_RESOURCE, id)
	return err
}