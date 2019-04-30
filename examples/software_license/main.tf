resource "onepassword_item_software_license" "this" {
  name  = "software-license"
  vault = "${var.vault_id}"

  main = {
    license_key = "${var.license_key}"
  }
}
