# onepassword_item_identity

This resource can load any identity from 1password.

## Example Usage

```hcl
data "onepassword_item_identity" "this" {
    name = "some-identity-from-vault"
}
```

## Argument Reference

* `name` - (Required) your identity title.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.
* `identification` - (Optional)
* `address` - (Optional)
* `internet` - (Optional)

The `identification` block support:

* `firstname` - (Optional)
* `initial` - (Optional)
* `lastname` - (Optional)
* `sex` - (Optional)
* `birth_date` - (Optional)
* `occupation` - (Optional)
* `company` - (Optional)
* `department` - (Optional)
* `job_title` - (Optional)

The `address` block support:

* `address` - (Optional)
* `default_phone` - (Optional)
* `home_phone` - (Optional)
* `cell_phone` - (Optional)
* `business_phone` - (Optional)

The `internet` block support:

* `username` - (Optional)
* `email` - (Optional)

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - identity id.
