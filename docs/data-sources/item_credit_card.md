# onepassword_item_credit_card

## Example Usage

```hcl
data "onepassword_item_credit_card" "this" {
    name = "credit_card-from-vault"
}
```

## Argument Reference

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

## Attribute Reference

In addition to the above arguments, the following attributes are exported:

* `id` - credit card id.
