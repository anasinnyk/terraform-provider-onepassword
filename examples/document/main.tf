resource "onepassword_item_document" "this" {
  name      = var.new_document_name
  vault     = var.vault_id
  file_path = "${path.module}/test.txt"
}
