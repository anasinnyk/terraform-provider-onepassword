package onepassword

import (
	"encoding/json"
)

const (
	// GroupResource is 1Password's internal designator for Groups
	GroupResource = "group"

	// GroupStateActive indicates an Active Group
	GroupStateActive = "A"

	// GroupStateDeleted indicates a Deleted Group
	GroupStateDeleted = "D"
)

// Group represents a 1Password Group resource
type Group struct {
	UUID  string
	Name  string
	State string
}

// ReadGroup gets an existing 1Password Group
func (o *OnePassClient) ReadGroup(id string) (*Group, error) {
	group := &Group{}
	res, err := o.runCmd(opPasswordGet, GroupResource, id)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, group); err != nil {
		return nil, err
	}
	return group, nil
}

// CreateGroup creates a new 1Password Group
func (o *OnePassClient) CreateGroup(v *Group) (*Group, error) {
	args := []string{opPasswordCreate, GroupResource, v.Name}
	res, err := o.runCmd(args...)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, v); err != nil {
		return nil, err
	}
	return v, nil
}

// DeleteGroup deletes a 1Password Group
func (o *OnePassClient) DeleteGroup(id string) error {
	return o.Delete(GroupResource, id)
}
