# onepassword_item_software_license

This resource can create any software license for 1password.

## Example Usage

```hcl
resource "onepassword_item_software_license" "this" {
  name        = "software-license-title"
  vault       = var.vault_id
  license_key = "SOME-SECURE-SOWTWARE-LICENSE-KEY"
}
```

## Argument Reference

* `name` - (Required) your software license title.
* `license_key` - (Optional) store your license key here.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - software license id.
