# onepassword_item_common

This resource can create any other item without required fields like Database/Membership/Wireless Router/Driver License/Outdoor License/Passport/Email Account/Reward Program/Social Security Number/Bank Account/Server in your 1password account.

## Example Usage

```hcl
resource "onepassword_item_common" "this" {
  name  = "Coupone"
  vault = "${var.vault_id}"

  template = "Reward Program"

  section = {
    field = {
      name   = "company name"
      string = "MacPaw"
    }

    field = {
      name   = "member name"
      string = "anasinnyk"
    }

    field = {
      name   = "member ID"
      string = "123"
    }

    field = {
      name      = "PIN"
      concealed = "123456qQ"
    }
  }

  section = {
    name = "More Information"

    field = {
      name   = "member ID (additional)"
      string = "321"
    }

    field = {
      name  = "customer service phone"
      phone = "+38 (000) 000 0000"
    }

    field = {
      name  = "phone for reservations"
      phone = "+38 (000) 000 0000"
    }

    field = {
      name = "website"
      url  = "https://groupon.com"
    }

    field = {
      name       = "member since"
      month_year = 201903
    }
  }
}
```

## Argument Reference

* `name` - (Required) your item title.
* `template` - (Required) your item category. Can be one of the next value `Database`, `Membership`, `Wireless Router`, `Driver License`, `Outdoor License`, `Passport`, `Email Account`, `Reward Program`, `Social Security Number`, `Bank Account`, `Server`.
* `vault` - (Optional) link to your vault, can be id (recommended) or name. If it's empty, it creates to default vault.
* `notes` - (Optional) note for this item.
* `tags` - (Optional) array of strings with any tag, for grouping your 1password item.
* `section` - (Optional) it's a block with additional information available in any other item type.

The `section` block support:

* `name` - (Optional) section title.
* `field` - (Optional) field in section.

The `field` block support:

* `name` - (Optional) field title.
* `string` - (Optional) if you have a text field use string.
* `url` - (Optional) if you have a URL field type (checks if URL is correct).
* `phone` - (Optional) if you have a phone number filed type.
* `email` - (Optional) if you have a email field type.
* `date` - (Optional) if you have a date field type should use a UNIXTIME.
* `month_year` - (Optional) if you have a month year field type, credit card expiration for example, use 6 number in next format `YYYYMM`.
* `totp` - (Optional) if you have a one time password you can save url in this type and 1password client can generate totp for you.
* `concealed` - (Optional) if you have a sensitive infromation, you can save it in this field type, it looks like a password.
* `sex` - (Optional) text field with information about geander, possible next vaules `male`,`female`.
* `card_type` - (Optional) text field with information about credit card type, possible next vaules `mc`, `visa`, `amex`, `diners`, `carteblanche`, `discover`, `jcb`, `maestro`, `visaelectron`, `laser`, `unionpay`.
* `reference` - (Optional) not supported yet. Potentially we can store reference between different items.
* `address` - (Optional) it's a address block.

*Note: MUST be one of there `string`,`url`,`phone`,`email`,`date`,`month_year`,`totp`,`concealed`,`address`,`sex`,`card_type`,`reference`.*

The `address` block support:

* `street` - (Optional) street information.
* `counrty` - (Optional) ISO2 country code.
* `state` - (Optional) state name.
* `region` - (Optional) region name.
* `city` - (Optional) city name.
* `zip` - (Optional) zip code.

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - item id.
