# onepassword_user

This resource can read user data in your 1Password account.

## Example Usage

```hcl
data "onepassword_user" "this" {
    email = "example@example.com"
}
```

## Argument Reference

* `email` - (Required) user email address.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - user id.
* `firstname` - User first name.
* `lastname` - User last name.
* `state` - Current user state. "A" for Active, "S" for Suspended.
