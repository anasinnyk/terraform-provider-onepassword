package onepassword

import (
	"encoding/json"
)

const (
	// UserResource is 1Password's internal designator for Users
	UserResource = "user"

	// UserStateActive indicates an Active User
	UserStateActive = "A"

	// UserStateSuspended indicates a Suspended User
	UserStateSuspended = "S"
)

// User represents a 1Password User resource
type User struct {
	UUID      string
	Email     string
	FirstName string
	LastName  string
	State     string
}

// ReadUser gets an existing 1Password User
// This supports multiple id parameter values, including "First Last", "Email", and "UUID".
func (o *OnePassClient) ReadUser(id string) (*User, error) {
	user := &User{}
	res, err := o.runCmd(opPasswordGet, UserResource, id)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(res, user); err != nil {
		return nil, err
	}
	return user, nil
}
