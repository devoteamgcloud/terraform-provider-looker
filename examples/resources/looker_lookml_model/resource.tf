resource "looker_lookml_model" "lookml_a" {
  name = "look-name"
  project_name = "project-gaeu"
  allowed_db_connection_names = ["bq-looker-connection-project-gaeu"]
}