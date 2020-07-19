# onepassword_vault

This resource can create vaults in your 1password account.

## Example Usage

```hcl
resource "onepassword_vault" "this" {
    name = "new-vault"
}
```

## Argument Reference

* `name` - (Required) vault name.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - vault id.
