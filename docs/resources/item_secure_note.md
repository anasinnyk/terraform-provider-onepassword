# onepassword_item_secure_note

This resource can create/load any secure note for/from 1password.

## Example Usage

```hcl
resource "onepassword_item_secure_note" "this" {
  name  = "secure-note-title"
  notes = <<<TEXT
    some multi line
    secret
    text
  >>>
  vault = "${var.vault_id}"
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
