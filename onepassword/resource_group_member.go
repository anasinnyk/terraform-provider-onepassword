package onepassword

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceGroupMemberRead,
		CreateContext: resourceGroupMemberCreate,
		DeleteContext: resourceGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"user": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
		},
	}
}

func resourceGroupMemberRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID, userID, err := resourceGroupMemberExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	v, err := m.onePassClient.ListGroupMembers(groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	var found string
	for _, member := range v {
		if member.UUID == userID {
			found = member.UUID
		}
	}

	if found == "" {
		d.SetId("")
		return nil
	}

	d.SetId(resourceGroupMemberBuildID(groupID, found))
	d.Set("group", groupID)
	d.Set("user", found)
	return nil
}

func resourceGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	m := meta.(*Meta)
	err := m.onePassClient.CreateGroupMember(
		d.Get("group").(string),
		d.Get("user").(string),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceGroupMemberBuildID(d.Get("group").(string), d.Get("user").(string)))
	return resourceGroupMemberRead(ctx, d, meta)
}

func resourceGroupMemberDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	groupID, userID, err := resourceGroupMemberExtractID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	m := meta.(*Meta)
	err = m.onePassClient.DeleteGroupMember(
		groupID,
		userID,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

// resourceGroupMemberBuildID will conjoin the group ID and user ID into a single string
// This is used as the resource ID.
//
// Note that user ID is being lowercased. Some operations require this user ID to be uppercased.
// Use the resourceGroupMemberExtractID function to correctly reverse this encoding.
func resourceGroupMemberBuildID(groupID, userID string) string {
	return strings.ToLower(groupID + "-" + strings.ToLower(userID))
}

// resourceGroupMemberExtractID will split the group ID and user ID from a given resource ID
//
// Note that user ID is being uppercased. Some operations require this user ID to be uppercased.
func resourceGroupMemberExtractID(id string) (groupID, userID string, err error) {
	spl := strings.Split(id, "-")
	if len(spl) != 2 {
		return "", "", fmt.Errorf("Improperly formatted group member string. The format \"groupid-userid\" is expected")
	}
	return spl[0], strings.ToUpper(spl[1]), nil
}
