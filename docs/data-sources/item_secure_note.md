# onepassword_item_secure_note

This resource can load any secure note from 1password.

## Example Usage

```hcl
data "onepassword_item_secure_note" "this" {
    name = "some-secure-note-from-vault"
}
```

## Argument Reference

* `name` - (Required) your secure note title.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common (main field for this type).
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - secure note id.
