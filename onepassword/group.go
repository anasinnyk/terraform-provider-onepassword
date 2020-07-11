package onepassword

import (
	"encoding/json"
	"fmt"
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

// ListGroupMembers lists the existing Users in a given Group
func (o *OnePassClient) ListGroupMembers(id string) ([]User, error) {
	users := []User{}
	if id == "" {
		return users, fmt.Errorf("Must provide an identifier to list group members")
	}

	res, err := o.runCmd(opPasswordList, "users", "--"+GroupResource, id)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, &users); err != nil {
		return nil, err
	}
	return users, nil
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

// CreateGroupMember adds a User to a Group
func (o *OnePassClient) CreateGroupMember(groupID string, userID string) error {
	args := []string{opPasswordAdd, UserResource, userID, groupID}
	_, err := o.runCmd(args...)
	return err
}

// UpdateGroup updates an existing 1Password Group
func (o *OnePassClient) UpdateGroup(id string, v *Group) error {
	args := []string{opPasswordEdit, GroupResource, id, "--name=" + v.Name}
	_, err := o.runCmd(args...)
	return err
}

// DeleteGroup deletes a 1Password Group
func (o *OnePassClient) DeleteGroup(id string) error {
	return o.Delete(GroupResource, id)
}

// DeleteGroupMember removes a User from a Group
func (o *OnePassClient) DeleteGroupMember(groupID string, userID string) error {
	args := []string{opPasswordRemove, UserResource, userID, groupID}
	_, err := o.runCmd(args...)
	return err
}
