provider "op" {}

data "op_vault" "this" {
  name = "Personal"
}

data "op_item_login" "666" {
  name = "SomeLogin666"
  vault = "${data.op_vault.this.id}"
}
data "op_item_login" "999" {
  name = "SomeLogin999"
  vault = "${data.op_vault.this.id}"
}

data "op_item_document" "this" {
  name = "Test TXT"
  vault = "${data.op_vault.this.id}"
}

output "doc_name" {
  value = "${data.op_item_document.this.file_name}"
}
output "doc_content" {
  value = "${data.op_item_document.this.content}"
}

resource "op_item_document" "this3" {
  name = "Some doc"
  tags = [
    "tag for doc",
    "system",
  ]
  vault = "${data.op_vault.this.id}"
  file_path = "./test.txt" 
}
resource "random_string" "secret" {
  length = "32"
}

resource "op_item_login" "888" {
  name = "New_888"
  vault = "${data.op_vault.this.id}"
  username = "admin"
  password = "${random_string.secret.result}"
  url = "https://whoma.ai"
  tags = [
    "this",
    "is",
    "Tagssssss",
  ]
  notes = "**Some** _awesome_ ~Markdown~"
  section = {
    name = "First Section"

    field = {
      name = "Address"
      address = {
        city = "New York"
        zip = "01001"
        country = "us"
        street = "Manheten"
        region = "Central Park"
        state = "NY"
      }
    }

    field = {
      name = "Text"
      string = "text value"
    }

    field = {
      name = "password"
      concealed = "concealed"
    }

    field = {
      name = "Website Url"
      url = "https://google.com"
    }
  }

  section = {
    name = "Second Section"

    field = {
      name = "Email address"
      email = "andr@kapitan.dev"
    }

    field = {
      name = "Telephone"
      phone = "+38 (111) 111 1111"
    }

    field = {
      name = "Current Date"
      date = 1556190244
    }

    field = {
      name = "Next Month and Year"
      month_year = 201905
    }

    field = {
      name = "2FA"
      totp = "otpauth://totp/label?secret=super-2fa-secret\u0026issuer=AWS"
    }
  }
}