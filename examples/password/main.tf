resource "onepassword_item_password" "this" {
  name     = "password"
  password = "${var.password}"
  vault    = "${var.vault_id}"
}
