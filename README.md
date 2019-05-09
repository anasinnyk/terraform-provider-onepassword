# Terraform OnePassword Provider

[![GolangCI](https://golangci.com/badges/github.com/anasinnyk/terraform-provider-1password.svg)](https://golangci.com/r/github.com/anasinnyk/terraform-provider-1password)
[![Build Status](https://travis-ci.com/anasinnyk/terraform-provider-1password.svg?branch=master)](https://travis-ci.com/anasinnyk/terraform-provider-1password)

## Table of Contents

[Provider](#Provider)

[Resources](#Resources)
* [onepassword_vault](#onepassword_vault)
* [onepassword_item_common](#onepassword_item_common)
* [onepassword_item_document](#onepassword_item_document)
* [onepassword_item_identity](#onepassword_item_identity)
* [onepassword_item_login](#onepassword_item_login)
* [onepassword_item_password](#onepassword_item_password)
* [onepassword_item_secure_note](#onepassword_item_secure_note)
* [onepassword_item_software_license](#onepassword_item_software_license)
* [onepassword_item_credit_card](#onepassword_item_credit_card)

[Data Sources](#data-sources)
* [onepassword_vault](#onepassword_vault-1)
* [onepassword_item_common](#onepassword_item_common-1)
* [onepassword_item_document](#onepassword_item_document-1)
* [onepassword_item_identity](#onepassword_item_identity-1)
* [onepassword_item_login](#onepassword_item_login-1)
* [onepassword_item_password](#onepassword_item_password-1)
* [onepassword_item_secure_note](#onepassword_item_secure_note-1)
* [onepassword_item_software_license](#onepassword_item_software_license-1)
* [onepassword_item_credit_card](#onepassword_item_credit_card-1)

## Provider

Terraform provider for 1password usage with your infrastructure, for example you can share password from your admin panel via some vault in you 1password company account

### Example Usage

```
provider "onepassword" {
    email      = "john.smith@example.com",
    password   = "super secret master password",
    secret_key = "A3-XXXXXX-XXXXXXX-XXXXX-XXXXX-XXXXX-XXXXX"
    sumdomain  = "company"
}
```

### Argument Reference

The following arguments are supported:

* `email` - (Required) your email address in 1password or via env variable `OP_EMAIL`.
* `password` - (Required) your master password from 1password or via env variable `OP_PASSWORD`.
* `secret_key` - (Required) secret key which you can download after registration or via env variable `OP_SECRET_KEY`.
* `subdomain` - (Optional) If you use corporate account you must fill subdomain form your 1password site. Defaults to `my` or via env variable `OP_SUBDOMAIN`.

## Resources

### onepassword_vault

### onepassword_item_common

### onepassword_item_document

### onepassword_item_identity

### onepassword_item_login

### onepassword_item_password

### onepassword_item_secure_note

### onepassword_item_software_license

### onepassword_item_credit_card

## Data Sources

### onepassword_vault

### onepassword_item_common

### onepassword_item_document

### onepassword_item_identity

### onepassword_item_login

### onepassword_item_password

### onepassword_item_secure_note

### onepassword_item_software_license

### onepassword_item_credit_card
