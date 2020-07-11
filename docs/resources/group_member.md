# onepassword_group_member

This resource can manage group membership within a 1Password group.

## Example Usage

### Resource

```hcl
resource "onepassword_group" "group" {
    group = "new-group"
}

data "onepassword_user" "user" {
    email = "example@example.com"
}

resource "onepassword_group_member" "example" {
    group = onepassword_group.group.id
    user = data.onepassword_user.user.id
}
```

## Argument Reference

* `group` - (Required) group id.
* `user` - (Required) user id.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - (Required) internal membership identifier.

## Import

1Password Group Members can be imported using the `id`, which consists of the group ID and user ID separated by a hyphen, e.g.

```
terraform import onepassword_group_member.example fmownretj6zdobn2cnjtqqyrae-KDLG56VTIJDXXBXC2KKCPHNHHI
```

**Note: this is case sensitive, and matches the case provided by 1Password.**
