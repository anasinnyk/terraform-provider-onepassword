# Terraform OnePassword Provider

[![GolangCI](https://golangci.com/badges/github.com/anasinnyk/terraform-provider-1password.svg)](https://golangci.com/r/github.com/anasinnyk/terraform-provider-1password)
[![Build Status](https://travis-ci.com/anasinnyk/terraform-provider-1password.svg?branch=master)](https://travis-ci.com/anasinnyk/terraform-provider-1password)

## Table of Contents

[Provider](#Provider)

* [onepassword_group](#onepassword_group)
* [onepassword_user](#onepassword_user)
* [onepassword_vault](#onepassword_vault)
* [onepassword_item_common](#onepassword_item_common)
* [onepassword_item_document](#onepassword_item_document)
* [onepassword_item_identity](#onepassword_item_identity)
* [onepassword_item_login](#onepassword_item_login)
* [onepassword_item_password](#onepassword_item_password)
* [onepassword_item_secure_note](#onepassword_item_secure_note)
* [onepassword_item_software_license](#onepassword_item_software_license)
* [onepassword_item_credit_card](#onepassword_item_credit_card)

## Provider

Terraform provider for 1password usage with your infrastructure, for example you can share password from your admin panel via some vault in you 1password company account. This provider based on 1Password CLI client version 0.5.5, but you can rewrite it by env variable `OP_VERSION`

### Example Usage

```hcl
provider "onepassword" {
    email      = "john.smith@example.com"
    password   = "super secret master password"
    secret_key = "A3-XXXXXX-XXXXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
    subdomain  = "company"
}
```

### Argument Reference

The following arguments are supported:

* `email` - (Optional) your email address in 1password or via env variable `OP_EMAIL`.
* `password` - (Optional) your master password from 1password or via env variable `OP_PASSWORD`.
* `secret_key` - (Optional) secret key which you can download after registration or via env variable `OP_SECRET_KEY`.
* `subdomain` - (Optional) If you use corporate account you must fill subdomain form your 1password site. Defaults to `my` or via env variable `OP_SUBDOMAIN`.

If `email`, `password` and `secret_key` is not set through the arguments or env variables, then the env variable `OP_SESSION_<subdomain>` is checked for existence. If set it will be assumed to be a valid session token and used while executing the `op` commands. Note that any dash `-` character within `subdomain` will be substituted upon `OP_SESSION_<subdomain>` env variable evaluation (e.g, if `subdomain=team-foo`, `OP_SESSION_team_foo` will be looked up).

## onepassword_group

This resource can create/load groups in your 1Password account.

### Example Usage

#### Resource

```hcl
resource "onepassword_group" "this" {
    name = "new-group"
}
```

#### Data Source

```hcl
data "onepassword_group" "this" {
    name = "exist-group"
}
```

### Argument Reference

* `name` - (Required) group name.

### Attribute Reference

* `id` - (Required) group id.

*Note: and all from arguments.*

### Import

1Password Groups can be imported using the `id`, e.g.

```
terraform import onepassword_group.group 7kalogoe3kirwf5aizotkbzrpq
```

## onepassword_user

This resource can read user data in your 1Password account.

### Example Usage

#### Data Source

```hcl
data "onepassword_user" "this" {
    email = "example@example.com"
}
```

### Argument Reference

* `email` - (Required) user email address.

### Attribute Reference

* `id` - (Required) user id.
* `firstname` - User first name.
* `lastname` - User last name.
* `state` - Current user state. "A" for Active, "S" for Suspended.

*Note: and all from arguments.*

## onepassword_vault

This resource create/load vault in your 1password account.

### Example Usage

#### Resource

```hcl
resource "onepassword_vault" "this" {
    name = "new-vault"
}
```

#### Data Source

```hcl
data "onepassword_vault" "this" {
    name = "exist-vault"
}
```

### Argument Reference

* `name` - (Required) vault name.

### Attribute Reference

* `id` - (Required) vault id.

*Note: and all from arguments.*

## onepassword_item_common

This resource can create/load any other item without required fields like Database/Membership/Wireless Router/Driver License/Outdoor License/Passport/Email Account/Reward Program/Social Security Number/Bank Account/Server in your 1password account.

### Example Usage

#### Resource

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

#### Data Source

```hcl
data "onepassword_item_common" "this" {
    name = "some-element-from-vault"
}
```

### Argument Reference

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

### Attribute Reference

* `id` - (Required) item id.

*Note: and all from arguments.*

## onepassword_item_document

This resource can create/load any document for/from 1password.

### Example Usage

#### Resource

```hcl
resource "onepassword_item_document" "this" {
  name      = "document-name"
  vault     = "${var.vault_id}"
  file_path = "${path.module}/test.txt"
}
```

#### Data Source

```hcl
data "onepassword_item_document" "this" {
    name = "some-document-from-vault"
}
```

### Argument Reference

* `name` - (Required) your document title.
* `field_path` - (Required) path to your document, which will be upload to 1password.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

### Attribute Reference

* `id` - (Required) document id.
* `content` - (Optional) document content.

*Note: and all from arguments.*

## onepassword_item_identity

This resource can create/load any identity for/from 1password.

### Example Usage

#### Resource

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

#### Data Source

```hcl
data "onepassword_item_identity" "this" {
    name = "some-identity-from-vault"
}
```

### Argument Reference

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

### Attribute Reference

* `id` - (Required) identity id.

*Note: and all from arguments.*

## onepassword_item_login

This resource can create/load any login for/from 1password.

### Example Usage

#### Resource

```hcl
resource "onepassword_item_login" "this" {
  name     = "login-title"
  username = "some-user-name"
  password = "123456qQ"
  url      = "https://example.com"
  vault    = "${var.vault_id}"
}
```

#### Data Source

```hcl
data "onepassword_item_login" "this" {
    name = "some-login-from-vault"
}
```

### Argument Reference

* `name` - (Required) your login title.
* `username` - (Optional) from this login.
* `password` - (Optional) from this login.
* `url` - (Optional) url for website from this login.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

### Attribute Reference

* `id` - (Required) login id.

*Note: and all from arguments.*

## onepassword_item_password

This resource can create/load any password for/from 1password.

### Example Usage

#### Resource

```hcl
resource "onepassword_item_password" "this" {
  name     = "login-title"
  password = "123456qQ"
  url      = "https://example.com"
  vault    = "${var.vault_id}"
}
```

#### Data Source

```hcl
data "onepassword_item_password" "this" {
    name = "some-password-from-vault"
}
```
### Argument Reference

* `name` - (Required) your password title.
* `password` - (Optional) store password here.
* `url` - (Optional) url for website from this password.
* `notes` - (Optional) see details in onepassword_item_common.
* `vault` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

### Attribute Reference

* `id` - (Required) password id.

*Note: and all from arguments.*

## onepassword_item_secure_note

This resource can create/load any secure note for/from 1password.

### Example Usage

#### Resource

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

#### Data Source

```hcl
data "onepassword_item_secure_note" "this" {
    name = "some-secure-note-from-vault"
}
```

### Argument Reference

* `name` - (Required) your secure note title.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common (main field for this type).
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

### Attribute Reference

* `id` - (Required) secure note id.

*Note: and all from arguments.*

## onepassword_item_software_license

This resource can create/load any software license for/from 1password.

### Example Usage

#### Resource

```hcl
resource "onepassword_item_software_license" "this" {
  name        = "software-license-title"
  vault       = "${var.vault_id}"
  license_key = "SOME-SECURE-SOWTWARE-LICENSE-KEY"
}
```

#### Data Source

```hcl
data "onepassword_item_software_license" "this" {
    name = "software-license-from-vault"
}
```

### Argument Reference

* `name` - (Required) your software license title.
* `license_key` - (Optional) store your license key here.
* `vault` - (Optional) see details in onepassword_item_common.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

### Attribute Reference

* `id` - (Required) software license id.

*Note: and all from arguments.*

## onepassword_item_credit_card

### Example Usage

#### Resource

```hcl
resource "onepassword_item_credit_card" "this" {
  name  = "Default Visa"
  vault = "${var.vault_id}"

  main = {
    cardholder  = "John Smith"
    type        = "visa"
    number      = "4111 1111 1111 1111"
    cvv         = "1111"
    expiry_date = 202205
    valid_from  = 201805
  }
}
```

#### Data Source

```hcl
data "onepassword_item_credit_card" "this" {
    name = "credit_card-from-vault"
}
```

### Argument Reference

* `name` - (Required) your credit card title.
* `vault` - (Optional) see details in onepassword_item_common.
* `main` - (Optional) block of card data.
* `notes` - (Optional) see details in onepassword_item_common.
* `tags` - (Optional) see details in onepassword_item_common.
* `section` - (Optional) see details in onepassword_item_common.

The `main` block support:

* `cardholder` - (Optional) store card holder name.
* `type` - (Optional) store card type value. see details in onepassword_item_common -> section -> field type card_type.
* `number` - (Optional) store 16 digit card numner.
* `cvv` - (Optional) sensitive data with your cvv card code.
* `expiry_date` - (Optional) store your exprite date in month year format. see details in onepassword_item_common -> section -> field type card_type
* `valid_from` - (Optional) store date when your card was publish in month year format. see details in onepassword_item_common -> section -> field type card_type

### Attribute Reference

* `id` - (Required) credit card id.

*Note: and all from arguments.*
