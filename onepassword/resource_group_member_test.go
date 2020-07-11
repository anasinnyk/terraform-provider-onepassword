package onepassword

import "testing"

func Test_resourceGroupMemberBuildID(t *testing.T) {
	want := "v3zk6wiptl42r7cmzbmf23unny-tgkw5a3cpbcu5end3lld3wckxi"
	got := resourceGroupMemberBuildID("v3zk6wiptl42r7cmzbmf23unny", "TGKW5A3CPBCU5END3LLD3WCKXI")

	if want != got {
		t.Error("Did not correctly conjoin the group and user IDs: " + got)
	}
}

func Test_resourceGroupMemberExtractID(t *testing.T) {
	wantGroup := "v3zk6wiptl42r7cmzbmf23unny"
	wantUser := "TGKW5A3CPBCU5END3LLD3WCKXI"
	gotGroup, gotUser, err := resourceGroupMemberExtractID("v3zk6wiptl42r7cmzbmf23unny-tgkw5a3cpbcu5end3lld3wckxi")

	if err != nil {
		t.Error(err)
	} else if wantGroup != gotGroup {
		t.Error("Did not correctly extract the group ID: " + gotGroup)
	} else if wantUser != gotUser {
		t.Error("Did not correctly extract the user ID: " + gotUser)
	}

	// Test malformed ID
	_, _, err = resourceGroupMemberExtractID("totally not the right id")
	if err == nil {
		t.Error("Error was not returned from malformed id")
	}
}
