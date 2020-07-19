# onepassword_item_identity

This resource can create any identity for 1password.

## Example Usage

```hcl
resource "onepassword_item_identity" "this" {
  name  = "Andrii Nasinnyk"
  vault = "${var.vault_id}"

  identification = {
    firstname  = "Andrii"
    initial    = "AN#24"
    lastname   = "Nasinnyk"
    sex        = "male"
    birth_date = 575553660         //unix-time
    occupation = "Play Basketball"
    company    = "HarshPhil"
    department = "Guards"
    job_title  = "Point Guard"
  }

  address = {
    address = {
      city    = "Kyiv"
      street  = "11 Line"
      country = "ua"
      zip     = "46000"
      region  = "Dniprovskii"
      state   = "Kyiv"
    }

    default_phone  = "+38 (000) 000 0000"
    home_phone     = "+38 (000) 000 0000"
    cell_phone     = "+38 (000) 000 0000"
    business_phone = "+38 (000) 000 0000"
  }

  internet = {
    username = "anasinnyk"
    email    = "andriy.nas@gmail.com"
  }
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
