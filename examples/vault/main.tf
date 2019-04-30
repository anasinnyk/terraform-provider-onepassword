provider "onepassword" {}

data "onepassword_vault" "this" {
  name = "Personal"
}

resource "onepassword_vault" "this" {
  name = "New"
}

# data "onepassword_item_login" "666" {
#   name = "SomeLogin666"
#   vault = "${data.onepassword_vault.this.id}"
# }
# data "onepassword_item_login" "999" {
#   name = "SomeLogin999"
#   vault = "${data.onepassword_vault.this.id}"
# }

# data "onepassword_item_document" "this" {
#   name = "Test TXT"
#   vault = "${data.onepassword_vault.this.id}"
# }

# output "doc_content" {
#   value = "${data.onepassword_item_document.this.content}"
# }

# resource "onepassword_item_document" "this3" {
#   name = "Some doc"
#   tags = [
#     "tag for doc",
#     "system",
#   ]
#   vault = "${data.onepassword_vault.this.id}"
#   file_path = "./test.txt" 
# }
# resource "random_string" "secret" {
#   length = "32"
# }

# data "onepassword_item_secure_note" "this" {
#   name = "RelatedNote"
# }

# resource "onepassword_item_secure_note" "sn" {
#   name = "New Secure Note"
#   notes = "${random_string.secret.result}"
#   tags = [
#     "this",
#     "is",
#     "Tagssssss",
#   ]
#   vault = "${data.onepassword_vault.this.id}" 
  
#   section = {
#     name = "First Section"

#     field = {
#       name = "Address"
#       address = {
#         city = "New York"
#         zip = "01001"
#         country = "us"
#         street = "Manheten"
#         region = "Central Park"
#         state = "NY"
#       }
#     }
#   }
# }

# resource "onepassword_item_common" "reward" {
#   name = "reward"
#   tags = [
#     "taggsss",
#     "new",
#   ]
#   vault = "${data.onepassword_vault.this.id}"
#   notes = "**Some** _awesome_ ~Markdown~"

#   template = "Reward Program"

#   section = {
#     field = {
#       name = "company name"
#       string = "MacPaw"
#     }
#     field = {
#       name = "member name"
#       string = "anasinnyk"
#     }
#   }

#   section = {
#     name = "More Information"
#     field = {
#       name = "member since"
#       month_year = 201903
#     }
#   }
# }


# resource "onepassword_item_software_license" "license" {
#   name = "license"
#   tags = [
#     "taggsss",
#     "new",
#   ]
#   vault = "${data.onepassword_vault.this.id}"
#   notes = "**Some** _awesome_ ~Markdown~"

#   main = {
#     title = "NON-EMPTY"

#     license_key = "ASDSADASDASDSADSADSDSSDSDSDADASDASDASDADDDSADSADD"
#   }

#   section = {
#     name = "extra"

#     field = {
#       name = "text"
#       string = "text"
#     }
#   }
# }

# resource "onepassword_item_identity" "8888888" {
#   name = "Andrii Indetity"
#   vault = "${data.onepassword_vault.this.id}"
#   tags = [
#     "this",
#     "is",
#     "Tagssssss",
#   ]
#   notes = "**Some** _awesome_ ~Markdown~"

#   identification = {
#     firstname = "Andrii"
#     initial = "AN"
#     lastname = "Nasinnyk"
#     sex = "male"
#     birth_date = 575553660 //unix-time
#     occupation = "occupation"
#     company = "company"
#     department = "department"
#     job_title = "job_title"
#     field = {
#       name = "extra_name"
#       string = "extra"
#     }
#   }

#   address = {
#     address = {
#       city = "Kyiv"
#       street = "11 Line"
#       country = "ua"
#       zip = "46000"
#       region = "Region"
#       state = "State"
#     }
#     default_phone = "+38000000000"
#     home_phone = "+38000000000"
#     cell_phone = "+38000000000"
#     business_phone = "+38000000000"
#     field = {
#       string = "extra"
#     }
#   }

#   internet = {
#     username = "nas1k"
#     email = "andriy.nas@gmail.com"
#     field = {
#       string = "extra"
#     }
#   }

#   section = {
#     name = "extra"

#     field = {
#       name = "text"
#       string = "text"
#     }
#   }
# }

# resource "onepassword_item_password" "8888888" {
#   name = "New_888_Password"
#   vault = "${data.onepassword_vault.this.id}"
#   password = "${random_string.secret.result}"
#   url = "https://whoma.ai"
#   tags = [
#     "this",
#     "is",
#     "Tagssssss",
#   ]
#   notes = "**Some** _awesome_ ~Markdown~"
#   section = {
#     name = "First Section"

#     field = {
#       name = "Address"
#       string = "text"
#     }
#   }
# }

# resource "onepassword_item_login" "8888888" {
#   name = "New_888_88888"
#   vault = "${data.onepassword_vault.this.id}"
#   username = "admin"
#   password = "${random_string.secret.result}"
#   url = "https://whoma.ai"
#   tags = [
#     "this",
#     "is",
#     "Tagssssss",
#   ]
#   notes = "**Some** _awesome_ ~Markdown~"
#   section = {
#     name = "First Section"

#     field = {
#       name = "Address"
#       address = {
#         city = "New York"
#         zip = "01001"
#         country = "us"
#         street = "Manheten"
#         region = "Central Park"
#         state = "NY"
#       }
#     }

#     field = {
#       name = "Text"
#       string = "text value"
#     }

#     field = {
#       name = "password"
#       concealed = "concealed"
#     }

#     field = {
#       name = "Website Url"
#       url = "https://google.com"
#     }
#   }

#   section = {
#     name = "Second Section"

#     field = {
#       name = "Email address"
#       email = "andr@kapitan.devS"
#     }

#     field = {
#       name = "Telephone"
#       phone = "+38 (111) 111 1111"
#     }

#     field = {
#       name = "Current Date"
#       date = 1556190244
#     }

#     field = {
#       name = "Next Month and Year"
#       month_year = 201905
#     }

#     field = {
#       name = "2FA"
#       totp = "otpauth://totp/label?secret=super-2fa-secret\u0026issuer=AWS"
#     }
#   }
# }

# resource "onepassword_item_credit_card" "cc_own_card" {
#   name = "CC Card Own"
#   vault = "${data.onepassword_vault.this.id}"
#   notes = "My UKRSIBBANK CARD"
#   tags = [
#     "cc",
#     "personal",
#   ]

#   main = {
#     title = "Rewrite Main Section Name"
#     cardholder = "Andrii Nasinnyk"
#     type = "mc"
#     number = "5351 0000 0000 0000"
#     cvv = "000"
#     expiry_date = 201905
#     valid_from = 201805
#     field = {
#       name = "extra"
#       concealed = "extra"
#     }
#   }

#   section = {
#     name = "extra section"
#     field = {
#       name = "text"
#       string = "text"
#     }
#   }
# }