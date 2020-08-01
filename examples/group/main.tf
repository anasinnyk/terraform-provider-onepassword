data "onepassword_group" "this" {
  name = var.exist_group_name
}

resource "onepassword_group" "this" {
  name = var.new_group_name
}
