resource "looker_permission_set" "set" {
  name = "My permission set"
  permissions = ["manage_homepage", "manage_spaces"]
}