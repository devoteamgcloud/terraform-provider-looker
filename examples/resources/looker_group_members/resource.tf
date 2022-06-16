data "looker_group" "group_two" {
  id = "4" // Parent group
}


resource "looker_group_member" "member_binding" {
  group_id = data.looker_group.group_two.id

  user {
    id = "1"
  }

  user {
    id = "3"
  }
}


resource "looker_group_member" "member_binding_secundo" {
  group_id = data.looker_group.group_two.id

  user {
    id = "4"
  }

  group {
    id = "3"
  }
}