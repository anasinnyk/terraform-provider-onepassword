---
layout: "onepassword"
page_title: "Provider: 1Password"
sidebar_current: "docs-onepassword-index"
description: |-
  The 1Password provider store your secret data in secure vaults.
---

# 1Password Provider

The 1Password provider store your secret data in secure vaults.

## Data Sources

* [Data Sources: onepassword_vault](vault.html)
* [Data Sources: onepassword_item_common](common.html)
* [Data Sources: onepassword_item_document](document.html)
* [Data Sources: onepassword_item_identity](identity.html)
* [Data Sources: onepassword_item_login](login.html)
* [Data Sources: onepassword_item_password](password.html)
* [Data Sources: onepassword_item_secure_note](secure_note.html)
* [Data Sources: onepassword_item_software_license](software_license.html)
* [Data Sources: onepassword_item_credit_card](credit_card.html)


## Resources

* [Resource: onepassword_vault](data-vault.html)
* [Resource: onepassword_item_common](data-common.html)
* [Resource: onepassword_item_document](data-document.html)
* [Resource: onepassword_item_identity](data-identity.html)
* [Resource: onepassword_item_login](data-login.html)
* [Resource: onepassword_item_password](data-password.html)
* [Resource: onepassword_item_secure_note](data-secure-note.html)
* [Resource: onepassword_item_software_license](data-software-license.html)
* [Resource: onepassword_item_credit_card](data-credit-card.html)

## Example Usage


## Authentication

There are generally two ways to configure the 1Password provider.

### Environment variables


### Statically defined credentials

## Argument Reference

The following arguments are supported:

* `email` - (Required) your email address in 1password.
* `password` - (Required) your master password from 1password.
* `secret_key` - (Required) secret key which you can download after registration.
* `subdomain` - (Optional) If you use corporate account you must fill subdomain form your 1password site. Defaults to `my`.