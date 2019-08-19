output "team" {
  value = data.onepassword_group.this.id
}

output "new" {
  value = onepassword_group.this.id
}
