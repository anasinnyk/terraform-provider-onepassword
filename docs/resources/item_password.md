# onepassword_item_password

This resource can create any password for 1password.

## Example Usage

```hcl
resource "onepassword_item_password" "this" {
  name     = "login-title"
  password = "123456qQ"
  url      = "https://example.com"
  vault    = var.vault_id
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
