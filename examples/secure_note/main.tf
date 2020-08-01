resource "onepassword_item_secure_note" "this" {
  name  = "secure_note"
  notes = var.secret
  vault = var.vault_id
}
