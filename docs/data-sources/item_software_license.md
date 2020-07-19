# onepassword_item_software_license

This resource can create/load any software license for/from 1password.

## Example Usage

```hcl
data "onepassword_item_software_license" "this" {
    name = "software-license-from-vault"
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
