package lookergo

const userAttributesSetBasePath = "4.0/user_attributes"

type UserAttributesResource interface {
}

type UserAttributesResourceOp struct {
	client *Client
}

var _ UserAttributesResource = &UserAttributesResourceOp{}

/*

List
Get
Create
Delete
GroupValueGet
GroupValueSet

### Get All User Attributes
GET {{endpoint}}/4.0/user_attributes


### Create User Attribute
POST {{endpoint}}/4.0/user_attributes


### Get User Attribute
GET {{endpoint}}/4.0/user_attributes/{{user_attribute_id}}


### Update User Attribute
PATCH {{endpoint}}/4.0/user_attributes/{{user_attribute_id}}


### Delete User Attribute
DELETE {{endpoint}}/4.0/user_attributes/{{user_attribute_id}}


### Get User Attribute Group Values
GET {{endpoint}}/4.0/user_attributes/{{user_attribute_id}}/group_values


### Set User Attribute Group Values
POST {{endpoint}}/4.0/user_attributes/{{user_attribute_id}}/group_values

*/
