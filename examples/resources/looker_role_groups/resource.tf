resource "looker_role_groups" "role_member" {
  role_id = "345"
  group {
    id = looker_group.group_admin.id
  }
}