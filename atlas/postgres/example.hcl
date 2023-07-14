schema "public" {
}

table "orders" {
  schema = schema.public

  column "id" {
    type = bigserial
  }
  column "no" {
    type = varchar
  }
  column "name" {
    type = varchar
  }
  column "description" {
    type = text
  }
  column "account" {
    type = varchar
  }
  column "customer" {
    type = varchar
  }
  column "email" {
    type = varchar
  }
  column "phone" {
    type = varchar
  }
  column "address" {
    type = varchar
  }
  column "price" {
    type = decimal
  }

  primary_key {
    columns = [
      column.id
    ]
  }
}
