resource "looker_role_member" "role_member" {
  target_role_id = "345"
  group {
    id = looker_group.group_admin.id
  }
}