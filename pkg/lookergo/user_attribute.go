package lookergo

import (
	"context"
	"fmt"
)

type UserAttributesResourceOp struct {
	client *Client
}

const UserAttributesBasePath = "4.0/user_attributes"

type UserAttribute struct {
	Can                        *map[string]bool `json:"can,omitempty"`                           // Operations the current user is able to perform on this object
	Id                         string           `json:"id,omitempty"`                            // Unique Id
	Name                       string           `json:"name"`                                    // Name of user attribute
	Label                      string           `json:"label"`                                   // Human-friendly label for user attribute
	Type                       string           `json:"type"`                                    // Type of user attribute ("string", "number", "datetime", "yesno", "zipcode")
	DefaultValue               string           `json:"default_value,omitempty"`                 // Default value for when no value is set on the user
	IsSystem                   *bool            `json:"is_system,omitempty"`                     // Attribute is a system default
	IsPermanent                *bool            `json:"is_permanent,omitempty"`                  // Attribute is permanent and cannot be deleted
	ValueIsHidden              *bool            `json:"value_is_hidden,omitempty"`               // If true, users will not be able to view values of this attribute
	UserCanView                *bool            `json:"user_can_view,omitempty"`                 // Non-admin users can see the values of their attributes and use them in filters
	UserCanEdit                *bool            `json:"user_can_edit,omitempty"`                 // Users can change the value of this attribute for themselves
	HiddenValueDomainWhitelist *string          `json:"hidden_value_domain_whitelist,omitempty"` // Destinations to which a hidden attribute may be sent. Once set, cannot be edited.
}

type UserAttributeGroupValue struct {
	Can             map[string]bool `json:"can,omitempty"`               // Operations the current user is able to perform on this object
	Id              string          `json:"id,omitempty"`                // Unique Id of this group-attribute relation
	GroupId         string          `json:"group_id,omitempty"`          // Id of group
	UserAttributeId string          `json:"user_attribute_id,omitempty"` // Id of user attribute
	ValueIsHidden   bool            `json:"value_is_hidden,omitempty"`   // If true, the "value" field will be null, because the attribute settings block access to this value
	Rank            int64           `json:"rank,omitempty"`              // Precedence for resolving value for user
	Value           string          `json:"value,omitempty"`             // Value of user attribute for group
}

type UserAttributesResource interface {
	Get(context.Context, int) (*UserAttribute, *Response, error)
	Create(context.Context, *UserAttribute) (*UserAttribute, *Response, error)
	Update(context.Context, string, *UserAttribute) (*UserAttribute, *Response, error)
	Delete(context.Context, string) (*Response, error)
	SetUserAttributeValue(context.Context, []UserAttributeGroupValue, string) (*[]UserAttributeGroupValue, *Response, error)
	GetUserAttributeValue(context.Context, string) (*[]UserAttributeGroupValue, *Response, error)
}

func (s *UserAttributesResourceOp) Get(ctx context.Context, UserAttributeId int) (*UserAttribute, *Response, error) {
	return doGetById(ctx, s.client, UserAttributesBasePath, UserAttributeId, new(UserAttribute))
}

func (s *UserAttributesResourceOp) Create(ctx context.Context, requestUserAttribute *UserAttribute) (*UserAttribute, *Response, error) {
	return doCreate(ctx, s.client, UserAttributesBasePath, requestUserAttribute, new(UserAttribute))
}

func (s *UserAttributesResourceOp) Update(ctx context.Context, UserAttributeId string, requestUserAttribute *UserAttribute) (*UserAttribute, *Response, error) {
	return doUpdate(ctx, s.client, UserAttributesBasePath, UserAttributeId, requestUserAttribute, new(UserAttribute))
}

func (s *UserAttributesResourceOp) Delete(ctx context.Context, UserAttributeId string) (*Response, error) {
	return doDelete(ctx, s.client, UserAttributesBasePath, UserAttributeId)
}

func (s *UserAttributesResourceOp) SetUserAttributeValue(ctx context.Context, userAtt []UserAttributeGroupValue, attributeId string) (*[]UserAttributeGroupValue, *Response, error) {
	path := fmt.Sprintf("%s/%s/group_values", UserAttributesBasePath, attributeId)
	return doAddValue(ctx, s.client, path, new([]UserAttributeGroupValue), userAtt)
}

func (s *UserAttributesResourceOp) GetUserAttributeValue(ctx context.Context, UserAttributeId string) (*[]UserAttributeGroupValue, *Response, error) {
	path := fmt.Sprintf("%s/%s/group_values", UserAttributesBasePath, UserAttributeId)
	return doGet(ctx, s.client, path, new([]UserAttributeGroupValue))
}
