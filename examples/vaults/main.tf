provider "op" {}

data "op_vault" "this" {
  name = "Personal"
}

data "op_item_login" "666" {
  name = "SomeLogin666"
  vault = "${data.op_vault.this.id}"
}

# resource "op_vault" "test2" {
#   name = "Test"
# }

# resource "op_vault" "import" {
#   name = "Personal"
#   users = []
# }

# resource "op_group" "this" {
#   name = "SomeGroup"
#   users = [
#       "email@here.com"
#   ]
# }

# resource "op_vault" "new" {
#   name = "New"
# }

# resource "op_item_login" "this" {
#   name     = "MyItem"
#   url      = "https://my.account.com"

#   username = "USERNAME"
#   password = "PASSWORD"

#   section = {
#     title = "first"

#     field = {
#       type = "URL"
#       title = "Main Website"
#       value = "http://example.com"
#     }

#     field = {
#       type = "address"
#       title = "Address"
#       value = {
#           street = "",
#           region = "",
#           country = "ua",
#           zipcode = "02160"
#       }
#     }
#   }

#   section = {
#     title = "second"

#     field = {
#       type = "concealed"
#       title = "Password"
#       value = "it's sensetive data"
#     }
#   }

#   vault = "${op_vault.new}"

#   tags = [
#     "some",
#     "my",
#     "tags",
#   ]
# }