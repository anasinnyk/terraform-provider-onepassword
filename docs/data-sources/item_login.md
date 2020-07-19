# onepassword_item_login

This resource can load any login from 1password.

## Example Usage

```hcl
data "onepassword_item_login" "this" {
    name = "some-login-from-vault"
}
```

## Argument Reference

* `name` - (Required) your login title.
* `username` - (Optional) from this login.
* `password` - (Optional) from this login.
* `url` - (Optional) url for website from this login.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - login id.
