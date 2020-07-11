package onepassword

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOnePassClient_ReadGroup(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		id string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		want            *Group
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `{ "uuid": "uniq", "name": "foo" }`, nil
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "get", "group", "uniq", "--session="},
			want:            &Group{UUID: "uniq", Name: "foo"},
		},
		{
			name: "bad json",
			fields: fields{
				runCmd: func() (string, error) {
					return `This was supposed to be JSON`, nil
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "get", "group", "uniq", "--session="},
			wantErr:         true,
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "get", "group", "uniq", "--session="},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			got, err := o.ReadGroup(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.ReadGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnePassClient.ReadGroup() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.ReadGroup() = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}

func TestOnePassClient_CreateGroup(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		v *Group
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		want            *Group
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `{ "uuid": "uniq", "name": "foo" }`, nil
				},
			},
			args:            args{v: &Group{Name: "foo"}},
			wantExecResults: []string{"op", "create", "group", "foo", "--session="},
			want:            &Group{UUID: "uniq", Name: "foo"},
		},
		{
			name: "bad json",
			fields: fields{
				runCmd: func() (string, error) {
					return `This was supposed to be JSON`, nil
				},
			},
			args:            args{v: &Group{Name: "foo"}},
			wantExecResults: []string{"op", "create", "group", "foo", "--session="},
			wantErr:         true,
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args:            args{v: &Group{Name: "foo"}},
			wantExecResults: []string{"op", "create", "group", "foo", "--session="},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			got, err := o.CreateGroup(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.CreateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnePassClient.CreateGroup() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.CreateGroup() = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}

func TestOnePassClient_UpdateGroup(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		id string
		v  *Group
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		want            *Group
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `{ "uuid": "uniq", "name": "foo" }`, nil
				},
			},
			args: args{
				id: "uniq",
				v:  &Group{Name: "foo"},
			},
			wantExecResults: []string{"op", "edit", "group", "uniq", "--name=foo", "--session="},
			want:            &Group{UUID: "uniq", Name: "foo"},
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args: args{
				id: "uniq",
				v:  &Group{Name: "foo"},
			},
			wantExecResults: []string{"op", "edit", "group", "uniq", "--name=foo", "--session="},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			err := o.UpdateGroup(tt.args.id, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.UpdateGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.UpdateGroup() = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}

func TestOnePassClient_DeleteGroup(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		id string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		want            *Group
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `{ "uuid": "uniq", "name": "foo" }`, nil
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "delete", "group", "uniq", "--session="},
			want:            &Group{UUID: "uniq", Name: "foo"},
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "delete", "group", "uniq", "--session="},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			err := o.DeleteGroup(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.DeleteGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.DeleteGroup() = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}

func TestOnePassClient_ListGroupMembers(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		id string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		want            []User
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `[ { "uuid": "uniq", "firstname": "Testy", "lastname": "Testerton" } ]`, nil
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "list", "users", "--group", "uniq", "--session="},
			want:            []User{{UUID: "uniq", FirstName: "Testy", LastName: "Testerton"}},
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args:            args{id: "uniq"},
			wantExecResults: []string{"op", "list", "users", "--group", "uniq", "--session="},
			wantErr:         true,
		},
		{
			name:    "error-missing-id",
			args:    args{id: ""},
			want:    []User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			got, err := o.ListGroupMembers(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.ListGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OnePassClient.ListGroupMembers() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.ListGroupMembers() exec = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}

func TestOnePassClient_CreateGroupMember(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		userID  string
		groupID string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `{ }`, nil
				},
			},
			args:            args{userID: "userName", groupID: "groupName"},
			wantExecResults: []string{"op", "add", "user", "groupName", "userName", "--session="},
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args:            args{userID: "userName", groupID: "groupName"},
			wantExecResults: []string{"op", "add", "user", "groupName", "userName", "--session="},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			err := o.CreateGroupMember(tt.args.userID, tt.args.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.ListGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.ListGroupMembers() exec = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}

func TestOnePassClient_DeleteGroupMember(t *testing.T) {
	type fields struct {
		runCmd func() (string, error)
	}
	type args struct {
		userID  string
		groupID string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantExecResults []string
		wantErr         bool
	}{
		{
			name: "success",
			fields: fields{
				runCmd: func() (string, error) {
					return `{ }`, nil
				},
			},
			args:            args{userID: "userName", groupID: "groupName"},
			wantExecResults: []string{"op", "remove", "user", "groupName", "userName", "--session="},
		},
		{
			name: "error",
			fields: fields{
				runCmd: func() (string, error) {
					return ``, fmt.Errorf("oops")
				},
			},
			args:            args{userID: "userName", groupID: "groupName"},
			wantExecResults: []string{"op", "remove", "user", "groupName", "userName", "--session="},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &mockOnePassConfig{
				runCmd: tt.fields.runCmd,
			}
			o := mockOnePassClient(config)

			err := o.DeleteGroupMember(tt.args.userID, tt.args.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnePassClient.ListGroupMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(config.execCommandResults, tt.wantExecResults) {
				t.Errorf("OnePassClient.ListGroupMembers() exec = %v, want %v", config.execCommandResults, tt.wantExecResults)
			}
		})
	}
}
