provider "op" {}

data "op_vault" "test" {
  name = "Semple"
}

resource "op_vault" "test2" {
  name = "Test"
}

resource "op_vault" "import" {
  name = "Personal"
  users = []
}

resource "op_group" "this" {
  name = "SomeGroup"
  users = [
      "email@here.com"
  ]
}

resource "op_vault" "new" {
  name = "New"
}

resource "op_item_login" "this" {
  name     = "MyItem"
  url      = "https://my.account.com"

  username = "USERNAME"
  password = "PASSWORD"

  section = {
    title = "first"

    field = {
      type = "URL"
      title = "Main Website"
      value = "http://example.com"
    }

    field = {
      type = "address"
      title = "Address"
      value = {
          street = "",
          region = "",
          country = "ua",
          zipcode = "02160"
      }
    }
  }

  section = {
    title = "second"

    field = {
      type = "concealed"
      title = "Password"
      value = "it's sensetive data"
    }
  }

  vault = "${op_vault.new}"

  tags = [
    "some",
    "my",
    "tags",
  ]
}