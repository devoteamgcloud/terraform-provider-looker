resource "looker_connection" "example" {
  name            = "test_connection"
  username        = "user"
  database        = "sf"
  dialect_name    = "snowflake"
  db_timezone     = "Europe/Amsterdam"
  port            = "443"
  host            = "db.snowflake.com"
  password        = "password"
  ssl             = true
  max_connections = 5
  tmp_db_name     = "temp_db"
  pdt_context_override {
    host = "db.snowflake.com"
    username = "user1"
    password = "password1"
    jdbc_additional_params = "database=my_db&warehouse=DEMO"
  }
}