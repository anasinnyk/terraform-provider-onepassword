resource "onepassword_item_common" "this" {
  name  = "Coupone"
  vault = "${var.vault_id}"

  template = "Reward Program"

  section {
    field {
      name   = "company name"
      string = "MacPaw"
    }

    field {
      name   = "member name"
      string = "anasinnyk"
    }

    field {
      name   = "member ID"
      string = "123"
    }

    field {
      name      = "PIN"
      concealed = "123456qQ"
    }
  }

  section {
    name = "More Information"

    field {
      name   = "member ID (additional)"
      string = "321"
    }

    field {
      name  = "customer service phone"
      phone = "+38 (000) 000 0000"
    }

    field {
      name  = "phone for reserva​tions"
      phone = "+38 (000) 000 0000"
    }

    field {
      name = "website"
      url  = "https://groupon.com"
    }

    field {
      name       = "member since"
      month_year = 201903
    }
  }
}
