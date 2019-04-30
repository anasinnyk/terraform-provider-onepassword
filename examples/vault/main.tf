data "onepassword_vault" "this" {
  name = "${var.exist_vault_name}"
}

resource "onepassword_vault" "this" {
  name = "${var.new_vault_name}"
}
