schema "public" {
}

table "users" {
  schema = schema.public
  column "id" {
    type = bigserial
  }
  column "name" {
    type = varchar
  }
  column "manager_id" {
    type = bigint
  }
  primary_key {
    columns = [
      column.id
    ]
  }
  index "idx_name" {
    columns = [
      column.name
    ]
    unique = true
  }
}
