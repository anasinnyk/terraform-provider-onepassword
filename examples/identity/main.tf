resource "onepassword_item_identity" "this" {
  name  = "Andrii Nasinnyk"
  vault = var.vault_id

  identification {
    firstname  = "Andrii"
    initial    = "AN#24"
    lastname   = "Nasinnyk"
    sex        = "male"
    birth_date = 575553660         //unix-time
    occupation = "Play Basketball"
    company    = "HarshPhil"
    department = "Guards"
    job_title  = "Point Guard"
  }

  address {
    address {
      city    = "Kyiv"
      street  = "11 Line"
      country = "ua"
      zip     = "46000"
      region  = "Dniprovskii"
      state   = "Kyiv"
    }

    default_phone  = "+38 (000) 000 0000"
    home_phone     = "+38 (000) 000 0000"
    cell_phone     = "+38 (000) 000 0000"
    business_phone = "+38 (000) 000 0000"
  }

  internet {
    username = "anasinnyk"
    email    = "andriy.nas@gmail.com"
  }
}
