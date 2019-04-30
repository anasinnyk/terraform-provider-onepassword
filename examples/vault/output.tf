output "personal" {
  value = "${data.onepassword_vault.this.id}"
}

output "new" {
  value = "${onepassword_vault.this.id}"
}
