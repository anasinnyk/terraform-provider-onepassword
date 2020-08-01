resource "onepassword_item_credit_card" "this" {
  name  = "Default Visa"
  vault = var.vault_id

  main {
    cardholder  = "John Smith"
    type        = "visa"
    number      = "4111 1111 1111 1111"
    cvv         = "1111"
    expiry_date = 202205
    valid_from  = 201805
  }
}
