# onepassword_group

This resource can load groups from your 1Password account.

## Example Usage

```hcl
data "onepassword_group" "this" {
    name = "exist-group"
}
```

## Argument Reference

* `name` - (Required) group name.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - group id.

## Import

1Password Groups can be imported using the `id`, e.g.

```
terraform import onepassword_group.group 7kalogoe3kirwf5aizotkbzrpq
```
