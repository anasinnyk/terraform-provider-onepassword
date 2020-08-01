provider "onepassword" {
  email      = "your@email.here"              // or use environment variable OP_EMAIL
  password   = "super-master-password-here"   // or use environment variable OP_PASSWORD
  secret_key = "secret-key-from-pdf-document" // or use environment variable OP_SECRET_KEY
  subdomain  = "company-domain"               // skip it or use my if you use personal 1password account or use environment variable OP_SUBDOMAIN
}

provider "random" {
  version = "~> 2.3"
}

terraform {
    required_version = ">= 0.12"
}

resource "random_string" "password" {
  length = "32"
}

module "group" {
  source = "./group"
}

module "user" {
  source = "./user"
  email = "example@example.com"
}

module "vault" {
  source = "./vault"
}

module "document" {
  source   = "./document"
  vault_id = module.vault.new
}

module "login" {
  source   = "./login"
  login    = "anasinnyk"
  password = random_string.password.result
  website  = "https://terraform.io"
  vault_id = module.vault.new
}

module "secret_note" {
  source   = "./secure_note"
  secret   = random_string.password.result
  vault_id = module.vault.new
}

module "password" {
  source   = "./password"
  password = random_string.password.result
  vault_id = module.vault.new
}

module "software_license" {
  source      = "./software_license"
  license_key = random_string.password.result
  vault_id    = module.vault.new
}

module "credit_card" {
  source   = "./credit_card"
  vault_id = module.vault.new
}

module "identity" {
  source   = "./identity"
  vault_id = module.vault.new
}

module "reward_program" {
  source   = "./reward_program"
  vault_id = module.vault.new
}
