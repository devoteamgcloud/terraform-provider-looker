resource "looker_role" "role" {
  name              = "admin"
  permission_set_id = ["1"]
  model_set_id      = looker_model_set.example.id
}