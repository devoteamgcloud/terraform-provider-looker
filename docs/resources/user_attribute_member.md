---
page_title: "looker_user_attribute_member Resource - terraform-provider-looker"
subcategory: ""
description: |-
  
---
# looker_user_attribute_member (Resource)

## Example Usage
```terraform
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
```

## Example Output
```terraform
% terraform show
# looker_group.my-group:
resource "looker_group" "my-group" {
    delete_on_destroy = true
    id                = "39"
    name              = "mygroup"
    parent_groups     = []
    roles             = []
}

# looker_user_attribute.att:
resource "looker_user_attribute" "att" {
    id              = "35"
    label           = "attribute label"
    name            = "attribute_name"
    type            = "string"
    user_can_edit   = true
    user_can_view   = true
    value_is_hidden = false
}

# looker_user_attribute_member.name:
resource "looker_user_attribute_member" "name" {
    id                = "-"
    user_attribute_id = "35"

    group {
        id    = "39"
        name  = "mygroup"
        value = "attribute-value"
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `user_attribute_id` (String)

### Optional

- `group` (Block Set) (see [below for nested schema](#nestedblock--group))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--group"></a>
### Nested Schema for `group`

Required:

- `value` (String)

Read-Only:

- `id` (String) The ID of this resource.
- `name` (String)
