# onepassword_item_login

This resource can create/load any login for/from 1password.

## Example Usage

```hcl
resource "onepassword_item_login" "this" {
  name     = "login-title"
  username = "some-user-name"
  password = "123456qQ"
  url      = "https://example.com"
  vault    = "${var.vault_id}"
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
