schema "example" {
  charset = "utf8mb4"
  collate = "utf8mb4_unicode_ci"
}

table "cities" {
  schema = schema.example
  column "id" {
    null           = false
    type           = bigint
    unsigned       = true
    auto_increment = true
  }
  column "name" {
    null = true
    type = varchar(255)
  }
  column "country_code" {
    null = true
    type = varchar(50)
  }
  column "state_code" {
    null = true
    type = varchar(50)
  }
  column "latitude" {
    null = true
    type = double
  }
  column "longitude" {
    null = true
    type = double
  }
  primary_key {
    columns = [column.id]
  }
  index "idx_cities_country_code" {
    columns = [column.country_code]
  }
  index "idx_cities_state_code" {
    columns = [column.state_code]
  }
}


table "users" {
  schema = schema.example
  column "id" {
    null = false
    type = int
  }
  column "username" {
    null = true
    type = varchar(200)
  }
  column "age" {
    null = true
    type = int
  }
  primary_key {
    columns = [column.id]
  }
}
