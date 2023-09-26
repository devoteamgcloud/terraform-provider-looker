resource "looker_connection" "example" {
  name            = "testdummy901"
  username        = "hegdgxme"
  database        = "hegdgxme"
  dialect_name    = "postgres"
  port            = "5432"
  host            = "surus.db.elephantsql.com"
  password        = "polite-sculpin"
  ssl             = true
  max_connections = 5
  user_attribute_fields = []
}