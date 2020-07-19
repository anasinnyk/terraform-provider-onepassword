# onepassword_item_password

This resource can create/load any password for/from 1password.

## Example Usage

```hcl
data "onepassword_item_password" "this" {
    name = "some-password-from-vault"
}
```
## Argument Reference

* `name` - (Required) your password title.
* `password` - (Optional) store password here.
* `url` - (Optional) url for website from this password.
* `notes` - (Optional) see details in onepassword_item_common.
* `vault` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - password id.
