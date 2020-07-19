# onepassword_vault

This resource create/load vault in your 1password account.

## Example Usage

```hcl
data "onepassword_vault" "this" {
    name = "exist-vault"
}
```

## Argument Reference

* `name` - (Required) vault name.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - vault id.
