resource "looker_user_attribute" "att" {
  name = "attribute_name"
  label = "attribute label"
  type = "string"
}

resource "looker_group" "my-group" {
  name = "mygroup"
}

resource "looker_user_attribute_member" "name" {
  user_attribute_id = looker_user_attribute.att.id
  group {
    id = looker_group.my-group.id
    value = "attribute-value"
  }
}