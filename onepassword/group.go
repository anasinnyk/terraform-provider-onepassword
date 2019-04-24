package onepassword

import (
	"encoding/json"
)

const GROUP_RESOURCE = "group"

type Group struct {
	Uuid string
	Name string
}

func (o *OnePassClient) ReadGroup(id string) (error, *Group) {
	Group := &Group{}
	err, res := o.runCmd(ONE_PASSWORD_COMMAND_GET, GROUP_RESOURCE, id)
	if err != nil {
		return err, nil
	}
	if err = json.Unmarshal(res, Group); err != nil {
		return err, nil
	}
	return nil, Group
}

func (o *OnePassClient) CreateGroup(v *Group) (error, *Group) {
	args := []string{ONE_PASSWORD_COMMAND_CREATE, GROUP_RESOURCE, v.Name}
	err, res := o.runCmd(args...)
	if err != nil {
		return err, nil
	}
	if err = json.Unmarshal(res, v); err != nil {
		return err, nil
	}
	return nil, v
}

func (o *OnePassClient) DeleteGroup(id string) error {
	err, _ := o.runCmd(ONE_PASSWORD_COMMAND_DELETE, GROUP_RESOURCE, id)
	return err
}