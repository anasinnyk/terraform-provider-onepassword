resource "onepassword_item_login" "this" {
  name     = var.login
  username = var.login
  password = var.password
  url      = var.website
  vault    = var.vault_id
}
