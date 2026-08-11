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

# Snowflake connection using key pair authentication
resource "looker_connection" "key_pair_example" {
  name               = "test_connection_key_pair"
  username           = "user"
  database           = "sf"
  dialect_name       = "snowflake"
  port               = "443"
  host               = "db.snowflake.com"
  uses_key_pair_auth = true
  certificate        = filebase64("rsa_key.p8")
  file_type          = ".p8"
  # For an encrypted private key, set password to the decryption passphrase
  # password         = "passphrase"
}